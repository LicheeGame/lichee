# lichee

## web

### minigame

1. wechat code2Session
2. rank

go build -gcflags "all=-N -l"



wx.request({
         url:'a.cpm',
         data:{ },
         header:{
             'content-type':'application/x-www-form-urlencoded',
             'Authorization':'Bearer  xxxxxxxxxxxxx',
         },
         method:'POST',
         success:function (res) {
             
         }
     })

  var header = {
    "Content-Type": "application/json",
    "X-Requested-With": 'XMLHttpRequest'
  };
  let token = wx.getStorageSync("token");
  if (token) {
    header["Authorization"] = `Bearer ${token}`;
  }