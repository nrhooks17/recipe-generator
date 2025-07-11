package main

import (
    "log"
    "strings"
    "image"
    "image/jpeg"
    "image/png"
    "os"
    "encoding/json"
    "encoding/base64"
    "bytes"
    "path/filepath"

    "recipe-generator/internal/api/config"
    "recipe-generator/internal/api/service"
    "recipe-generator/internal/api/handler"
)

func main() {
    //get config
    config, err := config.Load()

    if err != nil {
    	log.Fatalf("An error was encountered loading application config: %v\n", err)
    }

    //grab images
    base64Image := getBase64Image(config)
    
    // grab prompt. Toss this in an env variable when completed.
    aiPrompt := `I need the name, a description, prep time, cook time and servings made in addition to the ingredients. I also need in json format. The json should be an array of recipes with those fields stated above. In each recipe object, there should also be an ingredients array which is an array of objects, and a procedure step array which is just an array of strings. In the ingredients object, there should be unit of measurement field, unit amount field, and ingredient name. Any units of measurement that are fractions should be converted to decimal. The procedure array should be named "procedure_steps". Each procedure step should be exactly one action or at most two actions. Also for the json names, use snake case instead of camel case. Lastly, if the unit of measurement is not a "real" unit of measurement, i.e. a whole, or a pinch or handful or to taste, add an extra field to say so. Add a field called true_ingredient name. If the ingredient name has anything other than the ingredient name in the text, and that as the ingredient name and add the true ingredient name with just the ingredient name with no extra text. Also for each recipe, I need a combined confidence level of the accuracy of the optical character recognition. 0 being the least confident and 10 being the most. Do not infer any json fields if they don't exist. Put null if the field I asked for doesn't exist. Also, only return the json. Don't return other text.`
    
    // body of request
    body := constructBody(base64Image, aiPrompt)

    // Convert the map to JSON
    log.Println("Converting body to json")
    payloadByteArr, err := json.Marshal(body)
    log.Println("Finished converting body to json")	
    
    if err != nil {
    	log.Fatal(err)
    }
    
    payload := string(payloadByteArr)

    anthropicClient, err := service.InitializeAnthropicClient(config)

    if err != nil {
	log.Fatal(err)
    }

    responseJson, err := anthropicClient.Post(config, payload)

    if err != nil {
	log.Fatal(err)
    }

    log.Println("Response body: ", responseJson)

    // next I need to take take that json and save it to the database
}

// constructs the body of the request I want to send to anthropic.
func constructBody(base64Image string, aiPrompt string) map[string]interface{} {

    log.Printf("Constructing body of request.")

    body := map[string]interface{}{
    	"messages": []map[string]interface{}{
    	    {
    		"content": []map[string]interface{}{
    		    {
    			"type": "image",
    			"source": map[string]string{
    			    "type": "base64",
    			    "media_type": "image/jpeg",
    			    "data": base64Image,
    			},
    		    },
    		    {
    			"text": aiPrompt,
    			"type": "text",
    		    },
    		},
    		"role": "user",
    	    },
    	},
    	"max_tokens": 1024,
    	"model": "claude-3-7-sonnet-20250219",
    }
    
    log.Printf("Body construction complete.")

    return body
}

// function used to grab images from local directory. 
func getBase64Image(config *config.Config) string {
	//grab images
    imagePath := config.RecipeImagesLocation 

	// downsize image
    fileBytes, err := downsizeImageByBytes(imagePath, imagePath, 1000000)
    
    if err != nil {
    	log.Printf("Error downsizing image to \n %n", imagePath)
    	log.Fatal(err)
    }
    
    // convert byte array to base64
    log.Println("Converting image byte array to base64. ")
    base64Image :=  base64.StdEncoding.EncodeToString(fileBytes)
    log.Println("Finished converting image.")

    return base64Image
}

func downsizeImageByBytes(inputPath, outputPath string, maxBytes int) ([]byte, error) {
    // Open input file
    log.Printf("Opening input file: %s\n", inputPath)
    file, err := os.Open(inputPath)
    if err != nil {
    	return nil, err
    }
    defer file.Close()
    
    // Decode image
    log.Printf("Decoding image: %s\n", inputPath)
    img, format, err := image.Decode(file)
    if err != nil {
    	return nil, err
    }
    
    log.Printf("Finished decoding image: %s\n", inputPath)
    // Get original dimensions
    bounds := img.Bounds()
    width, height := bounds.Dx(), bounds.Dy()
    
    // Start with original size and reduce by 10% each iteration
    log.Printf("Downsizing image: %s\n", inputPath);	
    scale := 1.0
    for {
    	log.Printf("Downsizing image to scale: %f\n", scale);
    	newWidth := int(float64(width) * scale)
    	newHeight := int(float64(height) * scale)
    
    	// Create resized image
    	resized := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
    	for y := 0; y < newHeight; y++ {
    	    for x := 0; x < newWidth; x++ {
    		srcX := x * width / newWidth
    		srcY := y * height / newHeight
    		resized.Set(x, y, img.At(srcX, srcY))
    	    }
    	}
    
    	// Encode to buffer to check size
    	var buf bytes.Buffer
    	var encodeErr error
    	
    	if format == "png" || strings.ToLower(filepath.Ext(outputPath)) == ".png" {
    	    encodeErr = png.Encode(&buf, resized)
    	} else {
    	    encodeErr = jpeg.Encode(&buf, resized, &jpeg.Options{Quality: 90})
    	}
    	
    	if encodeErr != nil {
    	    return nil, encodeErr
    	}
    
    	//buf is the bytes of the resized image
    
    	// Check if size is acceptable
    	if buf.Len() <= maxBytes {
    	    // Write to output file
    	    return buf.Bytes(), nil 
    	}
    
    	// Reduce scale for next iteration
    	scale *= 0.9
    	if scale < 0.1 { // Prevent infinite loop
    	    return buf.Bytes(), nil
    	}
    }
}
