# 系统环境
ubuntu 16.04

go 1.13.4

docker CE 19.03.3

docker-compose 1.18

# 执行步骤

1. 要提前安装好go1.13, docker, docker-compose, IDE推荐goland, 最好用Linux环境...

2. 然后把项目拷到任意目录, 新版go不再硬性规定放在$GOPATH/src

    git clone https://github.com/KimRasak/SecKill.git

    cd SecKill

3. 下载依赖包

    go mod download 

4. 直接把抢购服务运行起来, 拉取镜像可能过程有点久
    
    docker-compose up seckill
    
    ![后台终端](./images/app_backend.jpg)
    
    或者测试一下能不能把mysql和redis单独运行起来
    
    docker-compose up -d mysql-service
    
    docker-compose up -d redis-service
    
5. 都成功后运行 docker ps 可以查看到三个服务容器的运行状态
    ![服务docker](./images/service_docker.jpg)

# 登录mysql查看用户和优惠券数据

1. mysql -udeveloper -p123456 -h 127.0.0.1

2. show databases; use dev; show tables;

   ![mysql](./images/mysql.jpg)
# 登录redis查看session等数据
1. redis-cli -h 127.0.0.1 -p 6379 -a 123456
  
   ![redis](./images/redis.jpg)