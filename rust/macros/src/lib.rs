use proc_macro::TokenStream;
use quote::quote;
use syn::{parse_macro_input, FnArg, ItemFn, Pat};

/// 日志属性宏
///
/// 用法示例:
/// ```
/// #[log]
/// fn add(a: i32, b: i32) -> i32 {
///     a + b
/// }
/// ```
#[proc_macro_attribute]
pub fn log(_attr: TokenStream, item: TokenStream) -> TokenStream {
    let mut func = parse_macro_input!(item as ItemFn);
    let func_ident = &func.sig.ident;

    // 提取参数名
    let param_names = func.sig.inputs.iter().filter_map(|arg| {
        if let FnArg::Typed(pat_type) = arg {
            if let Pat::Ident(pat_ident) = &*pat_type.pat {
                Some(&pat_ident.ident)
            } else {
                None
            }
        } else {
            None
        }
    }).collect::<Vec<_>>();

    // 生成参数名的字符串表示
    let param_names_str = param_names.iter().map(|ident| ident.to_string()).collect::<Vec<String>>().join(", ");
    _ = param_names_str;

    // 生成打印参数的代码
    let print_params = quote! {
        println!("macro: Function `{}` called with parameters", stringify!(#func_ident));
        #(
            println!("macro:  {} = {:?}", stringify!(#param_names), #param_names);
        )*
    };

    // 保存原始函数体
    let original_body = func.block.clone();

    // 修改函数体，添加日志逻辑
    let new_block = quote! {
        {
            #print_params
            let result = { #original_body };
            println!("macro: Function {} returned: {:?}", stringify!(#func_ident), result);
            result
        }
    };
    // Convert proc_macro2::TokenStream to proc_macro::TokenStream
    let new_block_ts: proc_macro::TokenStream = new_block.into();
    func.block = Box::new(parse_macro_input!(new_block_ts as syn::Block));

    let expanded = quote! {
        #func
    };

    expanded.into()
}

/// SQL 函数宏
///
/// 用法示例:
/// ```
/// let query = sql!("SELECT * FROM users WHERE id = 1");
/// ```
#[proc_macro]
pub fn sql(input: TokenStream) -> TokenStream {
    let sql_stmt = syn::parse_macro_input!(input as syn::LitStr).value();
    let expanded = quote! {
        {
            println!("macro: Executing SQL: {}", #sql_stmt);
            #sql_stmt
        }
    };

    expanded.into()
}

/// 简单的调试宏
///
/// 用法示例:
/// ```
/// debug!(x, y, z);
/// ```
#[proc_macro]
pub fn debug(input: TokenStream) -> TokenStream {
    let input = proc_macro2::TokenStream::from(input);
    let expanded = quote! {
        {
            println!("DEBUG: {} = {:?}", stringify!(#input), #input);
            #input
        }
    };

    expanded.into()
}
