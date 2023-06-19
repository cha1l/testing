### Contester

##Request example

 ```
 {
    "task_name": "Three sum",
    "language" : "cpp",
    "code": "#include <iostream>\nusing namespace std;\n\nint main() {\n    int a, b, c;\n    cin >> a >> b >> c;\n    cout << a + b + c  << endl;\n    return 0;\n}"
}
 ```
 
##Creating task for testing

```
{
    "name" : "Three sum",
    "text" : "You have three numbers. Print their sum",
    "duration" : 1000000000,
    "tests" : [
        {
            "input" : "2 3 1",
            "expected" : "6"
        },
        {
            "input" : "2 3 2",
            "expected" : "7"
        }
    ]
}
```
