package exponie

import (
	"context"
	"encoding/json"
	"errors"
	"exponie/pkg/fileutils"
	"github.com/rs/zerolog/log"
	"gopkg.in/h2non/gentleman.v2"
	"io"
	"os"
	"strings"
)

var client = gentleman.New().
	URL(WEB_HOSTNAME).
	SetHeader("X-Exponie-Client", "Wails")

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetDataset() []string {
	df, err := os.Open(fileutils.JoinConfigPath("words.txt"))
	if err != nil {
		log.Err(err).Msg("Failed to get dataset file.")
		return []string{}
	}
	defer fileutils.Close(df)
	b, err := io.ReadAll(df)
	if err != nil {
		log.Err(err).Msg("Failed to read dataset file.")
		return []string{}
	}
	s := string(b)
	return strings.Split(s, "\n")
}

func (a *App) GetDefinitions() map[string]string {
	definitions := make(map[string]string)

	df, err := os.Open(fileutils.JoinConfigPath("definitions.txt"))
	if err != nil {
		log.Err(err).Msg("Failed to get dataset file.")
		return definitions
	}
	defer fileutils.Close(df)
	b, err := io.ReadAll(df)
	if err != nil {
		log.Err(err).Msg("Failed to read dataset file.")
		return definitions
	}
	s := string(b)
	tokens := strings.Split(s, "\n")
	for _, token := range tokens {
		token := token

		if !strings.Contains(token, ": ") {
			log.Error().Int("line", len(definitions)).Str("text", token).Msg("Bad token, incomplete definition.")
			continue
		}

		items := strings.SplitN(token, ": ", 2)
		word, definition := items[0], items[1]
		definitions[word] = definition
	}
	return definitions
}

func (a *App) GetVersion() (float32, error) {
	var currentVersion float32 = -1.0

	versionFilePath := fileutils.JoinConfigPath("version.json")
	vf, err := os.Open(versionFilePath)
	if err != nil {
		return currentVersion, err
	}
	defer fileutils.Close(vf)
	bytes, err := io.ReadAll(vf)
	type version struct {
		Version float32 `json:"version"`
	}
	if err != nil {
		return currentVersion, err
	}
	var ver version
	if err = json.Unmarshal(bytes, &ver); err != nil {
		return currentVersion, err
	}
	currentVersion = ver.Version
	return currentVersion, nil
}

func (a *App) EnsureDataset() error {
	var currentVersion float32 = -1.0
	hasExistingVersion := true

	versionFilePath := fileutils.JoinConfigPath("version.json")
	vf, err := os.OpenFile(versionFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		} else {
			log.Warn().Str("file", versionFilePath).Msg("Cannot find version file.")
			if err = fileutils.MkdirParent(versionFilePath); err != nil {
				return err
			}
			hasExistingVersion = false
		}
	}
	type version struct {
		Version float32 `json:"version"`
	}
	if err == nil && vf != nil {
		defer fileutils.Close(vf)
		bytes, err := io.ReadAll(vf)
		if err != nil {
			return err
		}
		if len(bytes) != 0 {
			var ver version
			if err = json.Unmarshal(bytes, &ver); err != nil {
				return err
			}
			currentVersion = ver.Version
		} else {
			hasExistingVersion = false
		}
	}
	resp, err := client.Request().Method("GET").Path("/dataset/version.json").Send()
	if err != nil {
		return err
	}
	var ver version
	if err = resp.JSON(&ver); err != nil {
		return err
	}
	if ver.Version == currentVersion {
		return nil
	}
	resp, err = client.Request().Method("GET").Path("/dataset/words.txt").Send()
	if err != nil {
		return err
	}
	if err = resp.SaveToFile(fileutils.JoinConfigPath("words.txt")); err != nil {
		return err
	}
	resp, err = client.Request().Method("GET").Path("/dataset/definitions.txt").Send()
	if err != nil {
		return err
	}
	if err = resp.SaveToFile(fileutils.JoinConfigPath("definitions.txt")); err != nil {
		return err
	}
	txt, err := json.Marshal(ver)
	if err != nil {
		return err
	}
	if hasExistingVersion {
		if err = vf.Truncate(0); err != nil {
			return err
		}
		if _, err = vf.Seek(0, 0); err != nil {
			return err
		}
		if _, err = vf.Write(txt); err != nil {
			return err
		}
	} else {
		file, err := os.Create(versionFilePath)
		if err != nil {
			return err
		}
		defer fileutils.Close(file)
		if _, err = file.Write(txt); err != nil {
			return err
		}
	}
	return nil
}
