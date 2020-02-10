package app

var myapp App

func init()  {
	myapp = App{}
}

type App struct {
	server interface{}
	shutDownFunc func()
}

func SetShutDownFunc(f func()) {
	myapp.shutDownFunc = f
}

func GetShutDownFunc() func(){
	return myapp.shutDownFunc
}

func SetServer(s interface{}){
	myapp.server = s
}

func GetServer() interface{}{
	return myapp.server
}