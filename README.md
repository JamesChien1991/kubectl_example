# kubectl_example

## http server
```
{
"ORDER_NUM": "string"
"ENV_VAR": "string"
"ENV_VAR1": "string"
...
}
```
A http server處理build firmware request  
在k8s運行  
會將body內所有key, value當作環境變數傳遞給build firmware的container
## build firmware job
http server使用k8s的rbac token操作cluster  
建立一個job build firmware  
使用pvc掛載filestore，將成品上傳  
