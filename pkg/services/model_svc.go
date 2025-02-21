package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sagemakerruntime"

	"your-module-name/config"
)

var (
	cfg            = config.GetInstance()
	runtimeSmClient *sagemakerruntime.SageMakerRuntime
)

func init() {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(cfg.GetEnv("AWS_REGION", "us-west-2")),
		Credentials: credentials.NewStaticCredentials(cfg.GetEnv("AWS_ACCESS_KEY", ""), cfg.GetEnv("AWS_SECRET_KEY", ""), ""),
	})
	if err != nil {
		panic(err)
	}
	runtimeSmClient = sagemakerruntime.New(sess)
}

func queryEndpoint(endpointName string, queryJSON map[string]interface{}) (map[string]interface{}, error) {
	queryJSONFormatted, err := json.Marshal(queryJSON)
	if err != nil {
		return nil, err
	}

	input := &sagemakerruntime.InvokeEndpointInput{
		EndpointName: aws.String(endpointName),
		ContentType:  aws.String("application/json"),
		Body:         queryJSONFormatted,
	}

	result, err := runtimeSmClient.InvokeEndpoint(input)
	if err != nil {
		return nil, err
	}

	var response map[string]interface{}
	err = json.Unmarshal(result.Body, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func queryAILogoVectorizerService(queryJSON map[string]interface{}) (map[string]interface{}, error) {
	return queryEndpoint(cfg.GetEnv("LOGO_VECTORIZER_ENDPOINT", ""), queryJSON)
}

func queryLogoVectorizer(queryJSON map[string]interface{}) (map[string]interface{}, error) {
	body := map[string]interface{}{
		"model_name":  "logo_vectorizer_v1",
		"model_input": queryJSON,
	}
	return queryAILogoVectorizerService(body)
}

func QueryLogoVectorizerImage(imageB64 string) ([]float64, error) {
	queryJSON := map[string]interface{}{
		"type":   "images",
		"inputs": []string{imageB64},
	}

	result, err := queryLogoVectorizer(queryJSON)
	if err != nil {
		return nil, err
	}

	imageFeatures, ok := result["image_features"].([]interface{})
	if !ok || len(imageFeatures) == 0 {
		return nil, fmt.Errorf("invalid or empty image features")
	}

	features, ok := imageFeatures[0].([]float64)
	if !ok {
		return nil, fmt.Errorf("invalid image features format")
	}

	return features, nil
}

func main() {
	// Example usage
	imageB64 := "your_base64_encoded_image_here"
	features, err := QueryLogoVectorizerImage(imageB64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Image features:", features)
}
