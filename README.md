# BEMANI Finder
用于BEMANI的查分工具  
总之就是用于查分啦(  

IIDX部分是阿猫实现的  
SDVX部分是烧饼  

# 使用说明
http请求地址+端口+服务  
地址和端口在toml改  

## IIDX相关

例: http://localhost:9999/ (查看服务是否存活)  
例: http://localhost:9999/get?nick=摇滚大房子&max=6 (根据外号取歌名,max=最多取N个)  
例: http://localhost:9999/set?id=30053&nick=罪过的圣堂 (设置MID和外号,一个外号不能绑定多个MID，且MID必须存在)  
例: http://localhost:9999/del?&nick=罪过的圣堂 (删除一个外号)  
例: http://localhost:9999/nicks (查看当前服务器所有外号)  
例: http://localhost:9999/songs (查看当前服务器所有MID对应的歌名,从本地music_data.json读的)  
例: http://localhost:9999/reload (重新加载DB, 两个json，更新music_data.json时要用)  
  
## SDVX相关
例: http://localhost:9999/sdvx/get (获取sdvx所有曲目信息)  
例: http://localhost:9999/sdvx/get?id=999 (通过id获取曲目信息,注意: 不存在返回null)  
例: http://localhost:9999/sdvx/get?query=晕 (通过别名或者曲名匹配获取曲目信息,注意: 返回多个值)  
例: http://localhost:9999/sdvx/aliases (获取全部SDVX别名信息)  
例: http://localhost:9999/sdvx/aliases?id=693 (通过曲目id获取曲目别名)  
例: http://localhost:9999/sdvx/matchid?query=I (通过完全匹配名称获取到曲目id)  
例: http://localhost:9999/sdvx/matchid?query=i&isnocase=1 (通过完全匹配名称但是忽略大小写获取到曲目id)  
例: http://localhost:9999/sdvx/matchid?query=i&isnocase=1&isfuzzy=1 (模糊匹配所有曲名中包含"i"的曲目并且获取到id)  
例: http://localhost:9999/sdvx/matchid?query=晕船&isalias=1 (通过完全匹配别名获取到曲目id)  
例: http://localhost:9999/sdvx/matchid?query=BI&isnocase=1&isalias=1 (通过完全匹配别名但是忽略大小写获取到曲目id)  
例: http://localhost:9999/sdvx/matchid?query=晕船&isnocase=1&isalias=1&isfuzzy=1 (模糊匹配所有别名中包含"晕船"的曲目并且获取到id)  
例: http://localhost:9999/sdvx/existid?id=1394 (判断id是否存在)
例: http://localhost:9999/sdvx/addali?id=991&alias=test (给id为991的曲目添加test别名,"status": 0则是成功)  
例: http://localhost:9999/sdvx/delali?alias=test (删除别名test,"status": 0则是成功)  
例: http://localhost:9999/sdvx/reload (重新加载sdvx数据库, 更新music_db.xml或aliases.json时使用)  