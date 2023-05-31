package utils

import "fmt"

func PutBlob(connString string, container string, blobName, text string) (res *HttpPost) {

	credential, err := NewSharedKeyCredential(connString)
	if err != nil {
		fmt.Println("NewSharedKeyCredential err")
	}

	URI := "https://" + credential.AccountName() + ".blob.core.windows.net/" + container + "/" + blobName

	post := &HttpPost{
		URI:         URI,
		RequestBody: []byte(text),
	}
	var header map[string][]string
	header = make(map[string][]string)
	header["x-ms-blob-type"] = []string{"BlockBlob"}
	err = credential.HttpPutRequest(post, header)
	post.ResponseBody = []byte(XML2JSON(string(post.ResponseBody)))

	return post
}
