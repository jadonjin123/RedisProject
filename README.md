# Redis Project
Please read the entire README before doing anything!
This program reads and writes to Redis database using http requests. Here is what you need to install on your Linux machine first.
1. Install Golang
```
sudo apt update

sudo apt upgrade

sudo apt search golang-go

sudo apt search gccgo-go

sudo apt install golang-go
```

2. Install redis
```
sudo apt install lsb-release

curl -fsSL https://packages.redis.io/gpg | sudo gpg --dearmor -o /usr/share/keyrings/redis-archive-keyring.gpg

echo "deb [signed-by=/usr/share/keyrings/redis-archive-keyring.gpg] https://packages.redis.io/deb $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/redis.list

sudo apt-get update

sudo apt-get install redis
```

Here is how to use it
1. Navigate to the httpserver folder
2. Run this to build the project
```
go build main.go processor.go
```

3. Now run this command to start the server
```
./main.go
```
4. Now open another terminal and run curl commands to create, read, update, and delete (CRUD) entries in your database. Here are some examples
```
curl  -v -X POST 'http://localhost:3333/api/insert?myschedule=GetUp&myschedule=Eat&myschedule=Drink&foo=bar&foo=foobar'
```
```
curl  -v -X PUT 'http://localhost:3333/api/insert?myschedule=GetUp&myschedule=Eat&myschedule=Drink&foo=bar&foo=foobar'
```
```
curl -v  'http://localhost:3333/api/get?key1=myschedule&key2=foo'
```
Expected response:
{
myschedule: {
GetUp,
Eat,
Drink,
},
foo: {
bar,
foobar,
},
}
```
curl -v -X DELETE  'http://localhost:3333/api/delete?key1=myschedule&key2=foo'
```

Note that the response for GET will be in JSON format.
Each key you insert points to a list, so if you insert a key value pair where the key already exists, it will append the value to the list belonging to that key.
When you are inserting keys, you can also set a TTL value. The number you use will be in seconds. It will set a TTL for all of the keys included in the url. If you do not include a TTL field, there will not be a TTL for those keys. Example:
```
curl  -v -X POST 'http://localhost:3333/api/insert?myschedule=GetUp&ttl=100'
```
This will set the TTL for key "myschedule" to 100 seconds. You can set the ttl to -1 to remove the TTL.

### Warning: Running the tests will erase all key value pairs in the database
5. There are several tests that you can run. The Redis database must be empty for these tests to succeed. They are named TestGet, TestInsert, TestDelete, and TestNotFound. To run them, you can run:
```
go test -v -run TestName
```
You can also run all the tests at once by running 
```
go test
```
Note: In order for TestNotFound to succeed, the main executable must be running.
