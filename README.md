# LiFrame服务器架构

# windows环境下部署方式
# 数据库表创建
1.修改server\main\createtables.go 文件的数据库连接配置

2.执行create_db.bat

# 构建运行
1.执行build.bat

2.修改conf下各个服务器的数据库配置

3.执行run.bat

到此服务器启动成功，该默认方式会启动loginserver、gateserver、masterserver、worldserver、gameserver各一个服，

但是loginserver、gateserver、worldserver、gameserver支持分布式部署
可以通过修改run_by_conf.bat脚本中启动服务的配置，实现启动多个同一类型的服务

# 对应的demo客户端 https://github.com/llr104/LiFrameDemo

# linux环境下部署方式和windows雷同