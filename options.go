package yaspeech

// The Voice for the synthesized speech.
// You can choose one of the following voices:
// Female voice: alyss, jane, oksana and omazh.
// Male voice: zahar and ermil.
// Default value of the parameter: oksana.
type Voice string

// Voices of the synthesized speech.
const (
	VoiceAlyss  Voice = "alyss"
	VoiceJane   Voice = "jane"
	VoiceOksana Voice = "oksana"
	VoiceOmazh  Voice = "omazh"
	VoiceZahar  Voice = "zahar"
	VoiceErmil  Voice = "ermil"
)

// Emotion is emotional tone of the voice.
// Acceptable values:
// good — Cheerful and friendly.
// evil — Irritated.
// neutral (default) — Without emotion.
type Emotion string

//Emotional tone of the voice
const (
	EmotionGood    Emotion = "good"
	EmotionEvil    Emotion = "evil"
	EmotionNeutral Emotion = "neutral"
)

// The Format of the synthesized audio.
// Acceptable values:
// lpcm — Audio file is synthesized in the LPCM format with no WAV header. Audio characteristics:
// Sampling — 8, 16, or 48 kHz, depending on the sampleRateHertz parameter value.
// Bit depth — 16-bit.
// Byte order — Reversed (little-endian).
// Audio data is stored as signed integers.
// oggopus (default) — Data in the audio file is encoded using the OPUS audio codec and compressed using the OGG container format (OggOpus).
type Format string

// Formats of the synthesized audio.
const (
	FormatLpcm    Format = "lpcm"
	FormatOggOpus Format = "oggopus"
)

// Language of the synthesized speech
// Acceptable values:
// ru-RU (default) — Russian.
// en-US — English.
// tr-TR — Turkish.
type Language string

// Languages
const (
	LangRU Language = "ru-RU"
	LangEN Language = "en-US"
	LangTR Language = "tr-TR"
)

// SampleRateHertz is the sampling frequency of the synthesized audio.
// Used if format is set to lpcm. Acceptable values:
// 48000 (default) — Sampling rate of 48 kHz.
// 16000 — Sampling rate of 16 kHz.
// 8000 — Sampling rate of 8 kHz.1
type SampleRateHertz int

// The sampling frequency of the synthesized audio
const (
	SampleRate48000 SampleRateHertz = 48000
	SampleRate16000 SampleRateHertz = 16000
	SampleRate8000  SampleRateHertz = 8000
)

// Topic is the language model to be used for recognition.
// The closer the model is matched, the better the recognition result.
// You can only specify one model per request.
// Default parameter value: general.
type Topic string

const (
	// TopicGeneral - Short phrases containing 3-5 words on various topics, including search engine or website queries.
	TopicGeneral Topic = "general"
	// TopicMaps - Addresses and names of companies or geographical features.
	TopicMaps Topic = "maps"
	// TopicDates - Names of months, ordinal numbers, and cardinal numbers.
	TopicDates Topic = "dates"
	// TopicNames - First and last names and phone call requests.
	TopicNames Topic = "names"
	// TopicNumbers - Cardinal numbers from 1 to 999 and delimiters (dot, comma, and dash)
	TopicNumbers Topic = "numbers"
)

// // ProfanityFilter controls the profanity filter in recognized speech.
// // Acceptable values:
// // false (default) — Profanity is not excluded from recognition results.
// // true — Profanity is excluded from recognition results.
// type ProfanityFilter bool
