void callHelloWord(void *helloFunc) {
    void (*GetHelloWord)(void) = (void (*)(void))helloFunc;
    GetHelloWord();
}