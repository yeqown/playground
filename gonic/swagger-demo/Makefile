gen-doc:
	# 与 v3无异，只是不格式化生成的json文件
	GO111MODULE=off swagger generate spec -o swagger-demo.swagger.final.json -b ./mainC -c swagger-demo -m --compact

gen-doc-v3:
	# 支持扫描 swagger:model
	GO111MODULE=off swagger generate spec -o swagger-demo.swagger.v3.json -b ./mainC -c swagger-demo -m

gen-doc-v2:
	# 如果遇到找不到自己项目内部包的情况，且main包花里呼哨
	GO111MODULE=off swagger generate spec -o swagger-demo.swagger.v2.json -b ./mainC -c swagger-demo

gen-doc-v1:
	# main包花里胡哨，并不是一个单文件
	GO111MODULE=off swagger generate spec -o swagger-demo.swagger.v1.json -b ./mainC

gen-doc-v0:
	# 常规用法（适用于符合goswagger标准的项目）
	GO111MODULE=off swagger generate spec -o swagger-demo.swagger.v0.json

default: gen-doc-v0 gen-doc-v1 gen-doc-v2 gen-doc-v3 gen-doc 
