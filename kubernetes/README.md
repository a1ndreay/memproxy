Each node run a copy of memproxy which ensures the lack of external network knote.

#                       single node
#             /¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯\
#  user       memproxy   memcached       app     
#    | --------> |           |            |
#    |           | --------> |            |
#    |           | cache hit | cache miss |
#    | <---------+---------- | -------->  |
#    | <---------+-----------+----------- |