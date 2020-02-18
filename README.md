# LiFrame服务器架构，基于go开发 --武汉疫情期间宅在家中无所事事的成果

# 通过go命令获取项目 go get github.com/llr104/LiFrame

# windows环境下部署方式
# 数据库表创建
1.修改server\createtables.go 文件的数据库连接配置

2.执行create_db.bat

# 构建运行
1.执行build.bat

2.修改conf下各个服务器的数据库配置

3.执行run.bat

到此服务器启动成功，该默认方式会启动loginserver、gateserver、masterserver、worldserver、gameserver各一个服，

但是loginserver、gateserver、worldserver、gameserver支持分布式部署
可以通过修改run_by_conf.bat脚本中启动服务的配置，实现启动多个同一类型的服务



# linux环境下部署方式
1.修改server\createtables.go 文件的数据库连接配置

2.执行create_db.sh

# 构建运行
1.执行build.sh

2.修改conf下各个服务器的数据库配置

3.执行run.sh

# 对应的demo客户端 https://github.com/llr104/LiFrameDemo
# 众志成城，共同战疫，武汉加油，中国加油
# 如果你有对项目有建议或疑惑，QQ群1054084192欢迎拍砖

# 后面会基于该架构完成一个slg的三国游戏，slg的服务器代码会全部上传到GitHub，slg的客户端暂时不开源