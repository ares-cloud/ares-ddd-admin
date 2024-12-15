package token

import "testing"

const token = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJpbUNsb3VkIiwic3ViIjoie1widXNlcklkXCI6XCI5MjYzNzIxNTY2NTk1MDcyMFwiLFwicGxhdGZvcm1cIjpcIklPU1wiLFwidG9rZW5cIjpcIlwiLFwidGVuYW50SWRcIjpcIjEyMTYxMTgxNDQ0MDk2MFwiLFwiZXhwaXJlX3RpbWVfc2Vjb25kc1wiOjAsXCJyb2xlXCI6XCJVU0VSXCJ9IiwiZXhwIjoxNzAxNDU5MzgyLCJuYmYiOjE3MDEwOTkzODIsImlhdCI6MTcwMTA5OTM4Mn0.am9KPLmiMczMx2r5rOlbvv5m-MyXGCvvEJwWW4Phlao2eLPUiw7fqo47wD1X7t1_4HlD_DNHfqixL3SCOfKJIw"

func Test_verifyToken(t *testing.T) {
	tokener := Def()
	data := make(map[string]interface{})
	err := tokener.Verify(token, &data)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(data)
}
