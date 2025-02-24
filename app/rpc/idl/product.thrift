namespace go product

struct Product {
    1: i32 id;
    2: string name;
    3: double price;
    4: i32 stock;
    5: string image;
    6: string description;
    7: string created_at;
    8: bool is_deleted;
}

service ProductService {
    CreateProductResp createProduct(1: CreateProductReq req);
    UpdateProductResp updateProduct(1: UpdateProductReq req);
    RemoveProductResp removeProduct(1: RemoveProductReq req);
    GetProductByIDResp getProductByID(1: GetProductByIDReq req);
    ListProductResp listProduct(1: ListProductReq req);
    ListProductResp adminListProduct(1: ListProductReq req);
    SearchProductsResp searchProducts(1: SearchProductsReq req);
}

struct CreateProductReq {
    1: string name;
    2: string description;
    3: double price;
    4: i32 stock;
    5: string image;
}

struct CreateProductResp {
    1: Product product;
}

struct UpdateProductReq {
    1: i32 id;
    2: string name;
    3: string description;
    4: double price;
    5: i32 stock;
    6: string image;
}

struct UpdateProductResp {
    1: Product product;
}

struct RemoveProductReq {
    1: i32 id;
}

struct RemoveProductResp {
    1: bool success;
}

struct GetProductByIDReq {
    1: i32 id;
}

struct GetProductByIDResp {
    1: Product product;
}

struct SearchProductsReq {
    1: string query;
    2: i32 page_num;
    3: i32 page_size;
}

struct SearchProductsResp {
    1: list<Product> products;
    2: i64 total_count;
}

struct ListProductReq {
    1: i32 page_num;
    2: i32 page_size;
}

struct ListProductResp {
    1: list<Product> products;
    2: bool success;
    3: i64 total_count;
}