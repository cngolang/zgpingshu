#  中国评书网批量下载

###  编译
    go build .
    
#### 优化压缩编译
    go build -ldflags "-s -w"
###  压缩执行文件
upx -9 zgpingshu

#  使用说明

##  下载专辑 http://shantianfang.zgpingshu.com/575/#play   
    
    windows:  zgpingshu.exe  http://shantianfang.zgpingshu.com/575/#play   
    linux:  ./zgpingshu  http://shantianfang.zgpingshu.com/575/#play   

##  按顺序下载专辑的10个文件，已下载会自动忽略  
    
    windows:  zgpingshu.exe -q 10  http://shantianfang.zgpingshu.com/1040/#play   
    linux:  ./zgpingshu   -q 10  http://shantianfang.zgpingshu.com/1040/#play  

##  需重新下载某个节目音频，直接删除对应文件，重新执行命令即可，提示网络异常，直接ctrl+c，重新执行  
