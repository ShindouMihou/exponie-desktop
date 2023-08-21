package main

import (
	"bufio"
	"exponie/pkg/exponie"
	"exponie/pkg/fileutils"
	"exponie/pkg/utils"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/h2non/gentleman.v2"
	"gopkg.in/h2non/gentleman.v2/plugins/timeout"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	DEFINITIONS_API_HOSTNAME = "https://api.dictionaryapi.dev"
	DEFINITIONS_API_PATH     = "/api/v2/entries/en/"
)

const (
	DEFINITIONS_OUTPUT = "/dataset/definitions.txt"
)

var client = gentleman.New().
	URL(DEFINITIONS_API_HOSTNAME).
	Use(timeout.Request(5*time.Second)).
	SetHeader("X-Exponie-Client", "Wails")

type Word struct {
	Meanings []Meaning `json:"meanings"`
}

type Meaning struct {
	Definitions []Definition `json:"definitions"`
}

type Definition struct {
	Definition string `json:"definition"`
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	app := exponie.NewApp()
	args := os.Args[1:]
	if utils.Has(args, "--fetch") {
		if err := app.EnsureDataset(); err != nil {
			log.Panic().Err(err).Msg("Failed to pull dataset from origin.")
		}
	}
	dataset := app.GetDataset()
	if len(dataset) == 0 {
		log.Panic().Msg("There is no dataset available, try running with `--fetch` to fetch the dataset.")
	}
	chunkSize := len(dataset) / runtime.NumCPU()

	wg := sync.WaitGroup{}
	chunks := utils.Chunk(dataset, chunkSize)
	definitions := make(map[string]string)

	start := time.Now()

	wd, _ := os.Getwd()
	path := filepath.Join(wd, DEFINITIONS_OUTPUT)

	initial := ""

	existing, err := os.Open(path)
	if err == nil {
		defer fileutils.Close(existing)
		b, err := io.ReadAll(existing)
		if err != nil {
			log.Panic().Msg("Failed to read existing definitions.")
		}
		s := string(b)
		tokens := strings.Split(s, "\n")
		for _, token := range tokens {
			token := token

			if !strings.Contains(token, ": ") {
				log.Panic().Int("line", len(definitions)).Str("text", token).Msg("Bad token, incomplete definition.")
			}

			items := strings.SplitN(token, ": ", 2)
			word, definition := items[0], items[1]
			definitions[word] = definition
		}
		initial = s
		log.Info().Int("cursor", len(definitions)).Msg("Continuing after existing definitions.")
	}

	_ = fileutils.MkdirParent(path)
	file, err := os.Create(path)
	if err != nil {
		log.Panic().Msg("Failed to pre-create output file.")
	}
	defer fileutils.Close(file)
	writer := bufio.NewWriter(file)
	_, _ = writer.WriteString(initial)
	flushMutex, appendMutex := sync.Mutex{}, sync.Mutex{}

	ratelimited := make(chan bool, 1)
	var currentChunk int64

	log.Info().
		Int("chunk_size", chunkSize).
		Int("num_cpu", runtime.NumCPU()).
		Int("dataset_size", len(dataset)).
		Int("existing_definitions_size", len(definitions)).
		Msg("Running definition pulling with the following configuration.")

	// we are getting 24 less to account for potential 404s.
	if len(definitions) == (len(dataset) - 24) {
		log.Panic().Msg("All definitions have already been mapped.")
	}

	time.Sleep(5 * time.Second)
	for _, chunk := range chunks {
		atomic.AddInt64(&currentChunk, 1)
		chunk := chunk
		currentChunk := currentChunk

		wg.Add(1)
		go func() {
			defer wg.Done()
			select {
			case <-ratelimited:
				if flushMutex.TryLock() {
					log.Info().Msg("Flushing writer buffer after ratelimit.")
					if err := writer.Flush(); err != nil {
						log.Panic().Err(err).Msg("Failed to flush writer buffer.")
					}
					flushMutex.Unlock()
				}
				log.Error().Msg("Ratelimit encountered, waiting for 5 minutes.")
				time.Sleep(5 * time.Minute)
			default:
				progress := 0
				for _, word := range chunk {
					if _, ok := definitions[word]; ok {
						progress += 1
						continue
					}
					log.Info().
						Int64("chunk", currentChunk).
						Str("word", word).
						Str("progress", fmt.Sprint(progress, "/", len(chunk))).
						Msg("Requesting definition for the given word.")
					response, err := client.Request().Method("GET").Path(DEFINITIONS_API_PATH + word).Send()
					if err != nil {
						log.Err(err).
							Int64("chunk", currentChunk).
							Str("word", word).
							Msg("Failed to get a specific word.")
						continue
					}
					if response.StatusCode == http.StatusTooManyRequests {
						ratelimited <- true
						break
					}
					if !response.Ok {
						log.Error().
							Int("status_code", response.StatusCode).
							Int64("chunk", currentChunk).
							Str("word", word).
							Msg("Failed to get a specific word.")
						continue
					}
					var words []Word
					if err := response.JSON(&words); err != nil {
						log.Err(err).
							Int64("chunk", currentChunk).
							Str("word", word).
							Msg("Failed to unmarshal a specific word's response.")
						continue
					}
					if len(words) == 0 || len(words[0].Meanings) == 0 || len(words[0].Meanings[0].Definitions) == 0 {
						log.Err(err).
							Int64("chunk", currentChunk).
							Str("word", word).
							Msg("Word has no definitions.")
						continue
					}
					appendMutex.Lock()
					definitions[word] = words[0].Meanings[0].Definitions[0].Definition
					if _, err := writer.WriteString(word + ": " + definitions[word] + "\n"); err != nil {
						log.Err(err).
							Int64("chunk", currentChunk).
							Str("word", word).
							Msg("Failed to append to buffer writer.")
						continue
					}
					appendMutex.Unlock()
					progress += 1
					log.Info().
						Int64("chunk", currentChunk).
						Str("word", word).
						Str("progress", fmt.Sprint(progress, "/", len(chunk))).
						Msg("Completed definition for the word.")
				}
			}
		}()
	}

	wg.Wait()
	if err := writer.Flush(); err != nil {
		log.Panic().Err(err).Msg("Failed to flush writer buffer.")
	}
	log.Info().
		Str("time_taken", time.Since(start).String()).
		Msg("Completed definition set building.")
}
