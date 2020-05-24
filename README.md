# mildom_chat

## Usage

``` go
listener, err := GetListener(10467370)
if err != nil {
    return err
}
for msg := range listener {
    log.Println(msg.Username, msg.Text)
}
```
