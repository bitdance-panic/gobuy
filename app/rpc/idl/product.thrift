namespace go product

struct Product {
    1: i32 id;  // 来自 Base 的 ID
    2: string name;
    3: double price;
    4: i32 stock;
    5: string img;  // 对应 Go 结构体中的 Image
    6: string description;  // 对应 Go 结构体中的 Description
    7: string created_at;  // 来自 Base 的 CreatedAt，格式为字符串
    8: string updated_at;  // 来自 Base 的 UpdatedAt，格式为字符串
    9: bool is_deleted;  // 来自 Base 的 IsDeleted
}

service ProductService {
    CreateProductResponse createProduct(1: CreateProductRequest req);
    UpdateProductResponse updateProduct(1: UpdateProductRequest req);
    DeleteProductResponse deleteProduct(1: DeleteProductRequest req);
    GetProductByIDResponse getProductByID(1: GetProductByIDRequest req);
    SearchProductsResponse searchProducts(1: SearchProductsRequest req);
}

struct CreateProductRequest {
    1: string name;
    2: string description;
    3: double price;
    4: i32 stock;
    5: string img;
}

struct CreateProductResponse {
    1: Product product;
}

struct UpdateProductRequest {
    1: i32 id;
    2: string name;
    3: string description;
    4: double price;
    5: i32 stock;
    6: string img;
}

struct UpdateProductResponse {
    1: Product product;
}

struct DeleteProductRequest {
    1: i32 id;
}

struct DeleteProductResponse {
    1: bool success;
}

struct GetProductByIDRequest {
    1: i32 id;
}

struct GetProductByIDResponse {
    1: Product product;
}

struct SearchProductsRequest {
    1: string query;
}

struct SearchProductsResponse {
    1: list<Product> products;
}
