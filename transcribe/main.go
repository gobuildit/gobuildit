// The transcribe binary submits an audio sample to Google Cloud Platform's
// Speech-to-Text API for transcription.
package main

import (
	"fmt"
	"io/ioutil"
	"log"

	speech "cloud.google.com/go/speech/apiv1"
	"golang.org/x/net/context"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
)

func main() {
	ctx := context.Background()
	// speech.NewClient assumes an environment variable
	// GOOGLE_APPLICATION_CREDENTIALS that points to the service-account.json
	// file.
	client, err := speech.NewClient(ctx)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	// Read a local audio file into memory.
	data, err := ioutil.ReadFile("nhk-japanese.flac")
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}

	resp, err := client.Recognize(ctx, &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:        speechpb.RecognitionConfig_FLAC,
			SampleRateHertz: int32(16000), // Assume a sample rate of 16k Hz
			LanguageCode:    "ja-JP",      // Hard code Japanese as the sample's language.
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{Content: data},
		},
	})
	if err != nil {
		log.Fatalf("failed to recognize: %v", err)
	}

	// Print all results, in order of confidence of accurancy.
	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			fmt.Printf("\"%v\" (confidence=%3f)\n", alt.Transcript, alt.Confidence)
		}
	}
}
