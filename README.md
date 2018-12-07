#  中国评书网批量下载

    zgpingshu  linux 64版本
    zgpingshu.exe win10 64版本(不知win764是否正常工作)
    zgpingshu-win7_xp_32.exe win7/xp 32版本


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


使用了@WindGreen 兄弟写的多线程下载模块，https://github.com/WindGreen/fileload 感谢额
