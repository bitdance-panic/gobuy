namespace go user

// 定义请求和响应的结构体
//struct RegisterReq {
//    1: string email
//    2: string password
//    3: string confirm_password
//}
//
//struct RegisterResp {
//    1: i32 user_id
//}
struct RegisterReq {
    1: string email
    2: string password
    3: string confirm_password
    4: string username
}

struct RegisterResp {
    1: i32 user_id
    2: string message
}

struct LoginReq {
    1: string username
    2: string password
}

struct LoginResp {
    1: bool success
    2: i32 user_id
}

struct GetUserReq {
    1: string user_id
}


struct GetUserResp {
    1: bool success
    2: string user_id
    3: string email
    4: string username
}

struct User {
    1: i32 user_id
    2: string username
    3: string email
    4: string refresh_token
}

struct GetUsersReq {
    1: i32 page
    2: i32 page_size
}


struct GetUsersResp {
    1: list<User> users
    2: string message
}

struct UpdateUserReq {
    1: string user_id
    2: optional string email
    3: optional string username
}

struct UpdateUserResp {
    1: bool success
}

//软删除
struct DeleteUserReq {
    1: string user_id
}

struct DeleteUserResp {
    1: bool success
}

// 定义服务接口
service UserService {
    RegisterResp Register(1: RegisterReq req)
    LoginResp Login(1: LoginReq req)
    GetUserResp GetUser(1: GetUserReq req)
    UpdateUserResp UpdateUser(1: UpdateUserReq req)
    DeleteUserResp DeleteUser(1: DeleteUserReq req)
    GetUsersResp GetUsers(1: GetUsersReq req)
}