


```
def:  Response{data=Paged{list=[]model.DatasetBrief}}
path:       "/datasets": {
                "get": {
                        "responses": {
                            "200": {
                                "description": "OK",
                                "schema": {
                                    "allOf": [
                                        {
                                            "$ref": "#/definitions/api.Response"
                                        },
                                        {
                                            "type": "object",
                                            "properties": {
                                                "data": {
                                                    "allOf": [
                                                        {
                                                            "$ref": "#/definitions/api.Paged"
                                                        },
                                                        {
                                                            "type": "object",
                                                            "properties": {
                                                                "list": {
                                                                    "type": "array",
                                                                    "items": {
                                                                        "$ref": "#/definitions/model.DatasetBrief"
                                                                    }
                                                                }
                                                            }
                                                        }
                                                    ]
                                                }
                                            }
                                        }
                                    ]
                                }
                            }
                        }
                    }
            }
```

```
             "parameters": [
                    {
                        "description": "Upload file",
                        "name": "file",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
```


```

def:       // @Success 200 {string} pong       
path:       "get": {
                "summary": "ping test",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
```



```
                    {
                        "type": "string",
                        "name": "typeName",
                        "in": "query",           // query 参数
                        "required": true
                    },
                    {
                        "type": "integer",
                        "name": "id",
                        "in": "path",            // path 参数
                        "required": true
                    },    
                    {
                        "description": "email",
                        "name": "email",
                        "in": "body", ---         //body 参数
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }

```