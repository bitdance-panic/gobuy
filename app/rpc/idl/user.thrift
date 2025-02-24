namespace go user

struct RegisterReq {
    1: string email
    2: string password
    4: string username
}

struct RegisterResp {
    1: i32 user_id
    2: bool success
}

struct LoginReq {
    1: string email
    2: string password
}

struct LoginResp {
    1: bool success
    2: i32 user_id
}

struct GetUserReq {
    1: i32 user_id
}


struct GetUserResp {
    1: bool success
    2: i32 user_id
    3: string email
    4: string username
}

struct User {
    1: i32 user_id
    2: string username
    3: string email
    4: string refresh_token
}

struct AdminListUserReq {
    1: i32 page_num
    2: i32 page_size
}


struct AdminListUserResp {
    1: list<User> users
    2: string message
    3: i64 total_count
}

struct UpdateUserReq {
    1: i32 user_id
    2: optional string email
    3: optional string username
}

struct UpdateUserResp {
    1: bool success
}

//软删除
struct RemoveUserReq {
    1: i32 user_id
}

struct RemoveUserResp {
    1: bool success
}

struct BlockUserReq {
    1: string identifier 
    2: string reason
    3: i64 expires_at     
}

struct BlockUserResp {
    1: string block_id
    2: bool success
}

struct UnblockUserReq {
    1: string identifier     
}

struct UnblockUserResp {
    1: bool success
}


// 定义服务接口
service UserService {
    RegisterResp register(1: RegisterReq req)
    LoginResp login(1: LoginReq req)
    GetUserResp GetUser(1: GetUserReq req)
    UpdateUserResp updateUser(1: UpdateUserReq req)
    RemoveUserResp removeUser(1: RemoveUserReq req)
    BlockUserResp blockUser(1: BlockUserReq req)
    UnblockUserResp unblockUser(1: UnblockUserReq req)
    AdminListUserResp adminListUser(1: AdminListUserReq req)
}