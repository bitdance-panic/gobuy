namespace go agent

service AgentService {
    AskResp ask(1: AskReq req);
}


struct AskReq {
    1: string user_prompt;
    2: i32 user_id;
}

struct AskResp {
    1: string content;
}