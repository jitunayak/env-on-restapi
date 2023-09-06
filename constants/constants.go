package constants

var Title = `Sample Code Snippet For Postman Test 📝
------------------------------------------------------------------------`

// var Aws_url = `GET - http://localhost:8088/aws`
var Sample_code = `
	const {accessKeyId, secretKey, sessionToken} = pm.response.json();

	// for setting global level variables
	pm.globals.set("accessKeyId", accessKeyId);
	pm.globals.set("secretKey", secretKey);
	pm.globals.set("sessionToken", sessionToken);

	// or collection level variables
	pm.collectionVariables.set("accessKeyId", accessKeyId);
	pm.collectionVariables.set("secretKey", secretKey);
	pm.collectionVariables.set("sessionToken", sessionToken);

------------------------------------------------------------------------`
