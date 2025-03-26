/// 
/// Copyright (C) BABEC. All rights reserved.
/// 
/// SPDX-License-Identifier: Apache-2.0
/// 

/// ------user-----
/// fn init_contract(issue_limit, total_supply, manager_pk)
/// fn upgrade(issue_limit, total_supply, manager_pk)
/// fn register()
/// fn issue_amount(to, amount)
/// fn transfer(to, amount)
/// fn approve(spender, amount)
/// fn transfer_from(from, to, amount)
/// ------query-----
/// fn name()
/// fn symbol()
/// fn total_supply() -> total_supply
/// fn issued_amount() -> amount
/// fn balance_of() -> amount
/// fn allowance() -> amount
/// fn query_address() -> addr
/// fn query_other_address(pub_key)->addr
use contract_sdk_rust::sim_context;
use sim_context::*;

const VERSION: &str = "1.0.0";
const NAME: &str = "asset_management";
const SYMBOL: &str = "erc20";

/// get contract version
#[no_mangle]
pub extern "C" fn get_version() {
    let ctx = &mut sim_context::get_sim_context();
    ctx.ok(VERSION.as_bytes());
    ctx.log(VERSION);
}

/// get name
#[no_mangle]
pub extern "C" fn name() {
    let ctx = &mut sim_context::get_sim_context();
    ctx.ok(NAME.as_bytes());
}

/// get symbol
#[no_mangle]
pub extern "C" fn symbol() {
    let ctx = &mut sim_context::get_sim_context();
    ctx.ok(SYMBOL.as_bytes());
}

/// init
///
/// @param issue_limit 发行资产单笔限制, > 100
/// @param total_supply 发行资产总额, >1000
/// @param balance 合约创建人的初始资产
/// @param manager_pk 管理员公钥，以逗号','分割，不可超过5个。管理员具有分配资产（发钱）的权利。当无可用资产时，不可再发。
///
/// @ok creator address
#[no_mangle]
pub extern "C" fn init_contract() {
    sim_context::log("init start.");
    let ctx = &mut sim_context::get_sim_context();

    let pub_key_creator = ctx.get_creator_pub_key();
    let address = calc_address(&pub_key_creator);
    // get param
    let issue_limit_str = ctx.arg_as_utf8_str("issue_limit");
    let total_supply_str = ctx.arg_as_utf8_str("total_supply");
    let mut manager_pk_str = ctx.arg_as_utf8_str("manager_pk");

    // verify num
    let issue_limit = match parse_amount(&issue_limit_str) {
        Ok(val) => val,
        Err(msg) => {
            ctx.error(&msg);
            return;
        }
    };
    if issue_limit < 100 {
        ctx.error("issue_limit less 100.");
        return;
    }
    ctx.log(&("issue limit is ".to_string() + issue_limit_str.as_str()));
    let total_supply = match parse_amount(&total_supply_str) {
        Ok(val) => val,
        Err(msg) => {
            ctx.error(&msg);
            return;
        }
    };
    if total_supply < 1000 {
        ctx.error("total_supply less 1000.");
        return;
    }
    ctx.log(&("total supply is ".to_string() + total_supply_str.as_str()));

    // verify pk
    if manager_pk_str.len() == 0 {
        manager_pk_str = pub_key_creator;
    }
    let pks = manager_pk_str.split(",");
    let mut count = 0;
    for item in pks {
        if !is_pk(item) {
            ctx.error("manager_pk error");
            return;
        }
        count += 1;
        // register for manager
        let address_item = calc_address(&item); // hash160
        let r = ctx.get_state(&address_item, "balance");
        let balance_item = r.unwrap();
        if balance_item.len() == 0 {
            ctx.put_state(&address_item, "balance", "0".as_bytes());
        }
    }
    if count > 5 {
        ctx.error("manager_pk too much ");
        return;
    }
    ctx.log(&("manager pk count is ".to_string() + count.to_string().as_str()));

    let mut issued_amount = 0;
    let r = ctx.get_state("init", "issued_amount");
    if r.is_ok() {
        issued_amount = get_u32_from_str_byte(&r.unwrap());
    }
    if issued_amount > total_supply {
        ctx.error("issued_amount more than total_supply");
        return;
    }
    ctx.log(&("issued amount is ".to_string() + issued_amount.to_string().as_str()));

    // save data
    ctx.put_state("init", "total_supply", total_supply_str.as_bytes());
    ctx.put_state("init", "manager_pk", manager_pk_str.as_bytes());
    ctx.put_state("init", "issue_limit", issue_limit_str.as_bytes());
    ctx.put_state(
        "init",
        "issued_amount",
        issued_amount.to_string().as_bytes(),
    );
    // log
    ctx.log(&format!(
        "init: successed. Address[{:?}] . Total_supply is {}. Issue limit is {}",
        address, total_supply, issue_limit
    ));
    ctx.ok(address.as_bytes());
}

/// upgrade
#[no_mangle]
pub extern "C" fn upgrade() {
    sim_context::log("upgrade");
    init_contract();
}

/// register 开户
///
/// 开户后，余额为0。每个证书只能开户一次
/// @ok addr
#[no_mangle]
pub extern "C" fn register() {
    let ctx = &mut sim_context::get_sim_context();

    let pub_key = ctx.get_sender_pub_key();
    let address_a = calc_address(&pub_key); // hash160
    let r = ctx.get_state(&address_a, "balance");
    if r.is_err() {
        ctx.error("Register fail");
        return;
    }
    let balance_a = r.unwrap();
    if balance_a.len() != 0 {
        ctx.error("Already exists.");
        return;
    }
    ctx.put_state(&address_a, "balance", "0".as_bytes());
    ctx.ok(address_a.as_bytes());
    ctx.log(&format!("[{:?}] register success", address_a));
}

/// issue_amount 发钱
///
/// 只有管理员有权限发钱
///
/// @param to 收账人地址，开户时返回地址，或者可查询地址
/// @param amount 金额，正整数
#[no_mangle]
pub extern "C" fn issue_amount() {
    log("1");
    let ctx = &mut sim_context::get_sim_context();
    log("12");
    let sender_pub_key = ctx.get_sender_pub_key();
    log(&sender_pub_key);
    let manager_pub_key = get_manager_pk(ctx).unwrap();
    log(&manager_pub_key);
    if !manager_pub_key.contains(&sender_pub_key) {
        ctx.log("No permission.");
        ctx.error("No permission.");
        return;
    }
    log("2");
    // let creator_pub_key = ctx.get_creator_pub_key();
    // if creator_pub_key != sender_pub_key {
    //     ctx.error("No permission.");
    //     return;
    // }

    let issue_limit = get_issue_limit(ctx).unwrap();
    let address_to = ctx.arg_as_utf8_str("to");
    let amount = ctx.arg_as_utf8_str("amount");

    let msg = format!(
        "issue_limit {} address_to {:?} amount{:?}",
        issue_limit, address_to, amount
    );
    log(&msg);
    let amount = match parse_amount(&amount) {
        Ok(val) => val,
        Err(msg) => {
            ctx.log(&msg);
            ctx.error(&msg);
            return;
        }
    };
    if amount > issue_limit {
        ctx.error(&format!(
            "amount[{}] more than issue limit {}.",
            amount, issue_limit
        ));
        return;
    }

    let balance_to = match get_balance(ctx, &address_to) {
        Ok(val) => val,
        Err(msg) => {
            ctx.error(&msg);
            return;
        }
    };

    let balance_to_new = balance_to + amount;
    let supply = get_total_supply(ctx).unwrap();
    let mut issued_amount = get_issued_amount(ctx).unwrap();
    issued_amount = issued_amount + amount;
    if issued_amount > supply {
        ctx.error(&format!(
            "Sorry, fund pool exceeds upper limit. Supply is {}. issued_amount is {}",
            supply, issued_amount
        ));
        return;
    }

    ctx.put_state(
        &address_to,
        "balance",
        balance_to_new.to_string().as_bytes(),
    );
    ctx.put_state(
        "init",
        "issued_amount",
        issued_amount.to_string().as_bytes(),
    );
    ctx.log(&format!(
        "issue_amount: issue [{:?}] [{}], issued amount is [{}], supply is {}",
        address_to, amount, issued_amount, supply
    ));
}

/// transfer 转账
///
/// @param to 收账人地址，开户时返回地址，或者可查询地址
/// @param amount 金额，正整数
/// @ok "ok"
#[no_mangle]
pub extern "C" fn transfer() {
    let ctx = &mut sim_context::get_sim_context();

    let pub_key = ctx.get_sender_pub_key();
    let from = calc_address(&pub_key); // hash160
    let to = ctx.arg_as_utf8_str("to");
    let amount = ctx.arg_as_utf8_str("amount");

    if from == to {
        ctx.error("You can't transfer to yourself");
        return;
    }

    let amount = match parse_amount(&amount) {
        Ok(val) => val,
        Err(msg) => {
            ctx.error(&msg);
            return;
        }
    };

    let mut balance_from = match get_balance(ctx, &from) {
        Ok(val) => val,
        Err(msg) => {
            ctx.error(&msg);
            return;
        }
    };

    let mut balance_to = match get_balance(ctx, &to) {
        Ok(val) => val,
        Err(msg) => {
            ctx.error(&msg);
            return;
        }
    };
    // 余额不足
    if amount > balance_from {
        ctx.error("The balance of from not enough.");
        return;
    }
    balance_from = balance_from - amount;
    balance_to = balance_to + amount;
    ctx.put_state(&to, "balance", balance_to.to_string().as_bytes());
    ctx.put_state(&from, "balance", balance_from.to_string().as_bytes());
    ctx.log(&format!("transfer: [{:?}] to [{:?}] {}", from, to, amount));
    ctx.ok("ok".as_bytes());
}

/// approve 授权spender一定数额
///
/// @param  spender 花钱的人
/// @param  amount  数额
#[no_mangle]
pub extern "C" fn approve() {
    let ctx = &mut sim_context::get_sim_context();
    let address_from = ctx.get_sender_pub_key();
    let address_from = calc_address(&address_from);
    let address_spender = ctx.arg_as_utf8_str("spender");
    let amount = ctx.arg_as_utf8_str("amount");

    if !is_addr(&address_spender) {
        ctx.error("address error");
        return;
    }

    let amount = match parse_amount(&amount) {
        Ok(val) => val,
        Err(msg) => {
            ctx.error(&msg);
            return;
        }
    };

    let _ = match get_balance(ctx, &address_spender) {
        Ok(val) => val,
        Err(msg) => {
            ctx.error(&msg);
            return;
        }
    };

    ctx.put_state(
        &address_from,
        &address_spender,
        amount.to_string().as_bytes(),
    );
    ctx.log(&format!(
        "approve success: from {:?} spender {:?} amount {}",
        address_from, address_spender, amount
    ))
}

/// 代转账
///
/// @param from 付款人地址
/// @param to 收账人地址，开户时返回地址，或者可查询地址
/// @param amount 金额，正整数
/// @ok "ok"
#[no_mangle]
pub extern "C" fn transfer_from() {
    let ctx = &mut sim_context::get_sim_context();

    let pub_key = ctx.get_sender_pub_key();
    let address_my = calc_address(&pub_key);
    let from = ctx.arg_as_utf8_str("from");
    let to = ctx.arg_as_utf8_str("to");
    let amount = ctx.arg_as_utf8_str("amount");

    let mut balance_to = match get_balance(ctx, &to) {
        Ok(val) => val,
        Err(msg) => {
            ctx.error(&msg);
            return;
        }
    };
    let mut balance_from = match get_balance(ctx, &from) {
        Ok(val) => val,
        Err(msg) => {
            ctx.error(&msg);
            return;
        }
    };
    let amount = match parse_amount(&amount) {
        Ok(val) => val,
        Err(msg) => {
            ctx.error(&msg);
            return;
        }
    };
    let mut balance_allowance = match get_allowance(ctx, &from, &address_my) {
        Ok(val) => val,
        Err(_) => {
            ctx.error("You need to add allowance from_to");
            return;
        }
    };
    // 授权金额不足
    if amount > balance_allowance {
        ctx.error("The allowance of from_to not enough.");
        return;
    }
    // 余额不足
    if amount > balance_from {
        ctx.error("The balance of from not enough.");
        return;
    }

    balance_allowance = balance_allowance - amount;
    balance_to = balance_to + amount;
    balance_from = balance_from - amount;
    ctx.put_state(&from, &address_my, balance_allowance.to_string().as_bytes());
    ctx.put_state(&to, "balance", balance_to.to_string().as_bytes());
    ctx.put_state(&from, "balance", balance_from.to_string().as_bytes());
    ctx.log(&format!(
        "transferFrom: [{:?}] to [{:?}] {}",
        from, to, amount
    ));
    ctx.ok("ok".as_bytes());
}
/// total_supply 当前合约总金额
///
/// @ok amount
#[no_mangle]
pub extern "C" fn total_supply() {
    let ctx = &mut sim_context::get_sim_context();

    let r = ctx.get_state("init", "total_supply");
    let supply = r.unwrap();
    let supply = std::str::from_utf8(&supply).unwrap();
    ctx.ok(supply.as_bytes());
    ctx.log(&format!("total supply is [{:?}] ", supply));
}

/// issued_amount 当前合约已发金额
///
/// @ok amount
#[no_mangle]
pub extern "C" fn issued_amount() {
    let ctx = &mut sim_context::get_sim_context();

    let r = ctx.get_state("init", "issued_amount");
    let data = r.unwrap();
    let data = std::str::from_utf8(&data).unwrap();
    ctx.ok(data.as_bytes());
    ctx.log(&format!("total issued is [{:?}] ", data));
}

/// balance_of 获取自己账户余额
///
/// @param owner
/// @ok balance
#[no_mangle]
pub extern "C" fn balance_of() {
    let ctx = &mut sim_context::get_sim_context();

    let address = ctx.arg_as_utf8_str("owner");
    let balance = match get_balance(ctx, &address) {
        Ok(val) => val,
        Err(msg) => {
            ctx.error(&msg);
            return;
        }
    };

    ctx.ok(balance.to_string().as_bytes());
    ctx.log(&format!(
        "[{:?}] balance is [{:?}] ",
        address,
        balance.to_string()
    ));
}

/// 获取授权限额
///
/// @param spender 被授权人
/// @param owner 授权人
#[no_mangle]
pub extern "C" fn allowance() {
    let ctx = &mut sim_context::get_sim_context();
    ctx.log("allowance");

    let to = ctx.arg_as_utf8_str("spender");
    let from = ctx.arg_as_utf8_str("owner");
    let _ = match get_balance(ctx, &to) {
        Ok(val) => val,
        Err(msg) => {
            ctx.error(&msg);
            return;
        }
    };
    let _ = match get_balance(ctx, &from) {
        Ok(val) => val,
        Err(msg) => {
            ctx.error(&msg);
            return;
        }
    };
    let amount = match get_allowance(ctx, &from, &to) {
        Ok(val) => val,
        Err(msg) => {
            ctx.error(&msg);
            return;
        }
    };
    ctx.ok(amount.to_string().as_bytes());
}

/// query_address 查询自己的钱包地址
///
/// @ok  address
#[no_mangle]
pub extern "C" fn query_address() {
    let ctx = &mut sim_context::get_sim_context();
    ctx.log("query_address");

    let pub_key = ctx.get_sender_pub_key();
    ctx.log(&format!("[{:?}] pk", pub_key));
    let address = calc_address(&pub_key); // hash160
    ctx.log(&format!("[{:?}] address", address));
    ctx.ok(address.as_bytes());
    ctx.log(&format!("[{:?}] address is [{:?}]", pub_key, address));
}

/// query_address 查询他人的钱包地址
///
/// @param pub_key 公钥
/// @ok    address
#[no_mangle]
pub extern "C" fn query_other_address() {
    let ctx = &mut sim_context::get_sim_context();

    let pub_key = ctx.arg_as_utf8_str("pub_key");
    if !is_pk(&pub_key) {
        ctx.error("pub_key error.");
        return;
    }
    let address = calc_address(&pub_key);
    let _ = match get_balance(ctx, &address) {
        Ok(val) => val,
        Err(msg) => {
            ctx.error(&msg);
            return;
        }
    };
    ctx.ok(address.as_bytes());
    ctx.log(&format!("pub key [{:?}] address [{:?}]", pub_key, address));
}

/// 从链上获取余额并校验
///
/// ERR：查询失败、未注册
fn get_balance(ctx: &mut dyn SimContext, address: &str) -> Result<u32, String> {
    if address.len() == 0 {
        return Err("miss address".to_string());
    }
    let r = ctx.get_state(address, "balance");
    if r.is_err() {
        return Err(format!("[{:?}] get balance fail", address));
    }
    let balance = r.unwrap();
    if balance.len() == 0 {
        return Err(format!("[{:?}] not registered", address));
    }
    let balance = std::str::from_utf8(&balance).unwrap();
    let balance = balance.parse::<u32>().unwrap();
    Ok(balance)
}

/// 从链上获取单笔资产发行上限
fn get_issue_limit(ctx: &mut dyn SimContext) -> Result<u32, String> {
    let r = ctx.get_state("init", "issue_limit");
    if r.is_err() {
        return Err("get_issue_limit fail".to_string());
    }
    let data = r.unwrap();
    if data.len() == 0 {
        return Err("get_issue_limit fail".to_string());
    }
    let issue_limit_str = std::str::from_utf8(&data).unwrap();
    Ok(issue_limit_str.parse::<u32>().unwrap())
}

/// 从链上获取已发行资产
fn get_issued_amount(ctx: &mut dyn SimContext) -> Result<u32, String> {
    let r = ctx.get_state("init", "issued_amount");
    if r.is_err() {
        return Err("get_issued_amount fail".to_string());
    }
    let data = r.unwrap();
    if data.len() == 0 {
        return Err("get_issued_amount fail".to_string());
    }
    let issued_amount = std::str::from_utf8(&data).unwrap();
    Ok(issued_amount.parse::<u32>().unwrap())
}

/// 从链上获取资产总额
fn get_total_supply(ctx: &mut dyn SimContext) -> Result<u32, String> {
    let r = ctx.get_state("init", "total_supply");
    if r.is_err() {
        return Err("get_issued_amount fail".to_string());
    }
    let data = r.unwrap();
    if data.len() == 0 {
        return Err("get_total_supply fail".to_string());
    }
    let total_supply = std::str::from_utf8(&data).unwrap();
    Ok(total_supply.parse::<u32>().unwrap())
}

/// 获取管理员
fn get_manager_pk(ctx: &mut dyn SimContext) -> Result<String, String> {
    let r = ctx.get_state("init", "manager_pk");
    if r.is_err() {
        return Err("get_issued_amount fail".to_string());
    }
    let data = r.unwrap();
    if data.len() == 0 {
        return Err("get_total_supply fail".to_string());
    }
    let manager_pk = std::str::from_utf8(&data).unwrap();
    Ok(manager_pk.to_string())
}

/// 获取授权总额
fn get_allowance(ctx: &mut dyn SimContext, from: &str, to: &str) -> Result<u32, String> {
    if from.len() == 0 || to.len() == 0 {
        return Err("miss address".to_string());
    }
    let r = ctx.get_state(from, to);
    if r.is_err() {
        return Err("get_issued_amount fail".to_string());
    }
    let data = r.unwrap();
    if data.len() == 0 {
        return Err("get_total_supply fail".to_string());
    }
    if data.len() == 0 {
        return Err(format!("[{:?}] not approved {:?}", from, to));
    }
    let amount = std::str::from_utf8(&data).unwrap();
    Ok(amount.parse::<u32>().unwrap())
}

/// 校验是否是公钥格式
fn is_pk(pk: &str) -> bool {
    if pk.len() == 0 || pk.len() > 64 {
        return false;
    }
    // for elem in pk.chars() {
    //     if '0' <= elem && '9' >= elem {
    //         continue;
    //     }
    //     if 'a' <= elem && 'z' >= elem {
    //         continue;
    //     }
    //     if 'A' <= elem && 'Z' >= elem {
    //         continue;
    //     }
    //     return false;
    // }
    true
}
/// 校验是否是公钥格式
fn is_addr(addr: &str) -> bool {
    if addr.len() == 0 || addr.len() > 64 {
        return false;
    }
    return true;
}

/// 解析金额为正整数
///
/// 为空、不是数字，非整数，超过int 则返回错误信息
fn parse_amount(amount: &str) -> Result<u32, String> {
    if amount.len() == 0 {
        return Err("amount is null".to_string());
    }
    for elem in amount.chars() {
        if '0' <= elem && '9' >= elem {
            continue;
        }
        if '-' == elem {
            return Err(format!("amount[{}] less than 0. ", amount));
        }
        return Err(format!("amount[{}] is  not number. ", amount));
    }
    // let amount_reg = Regex::new(r"^\d+$").unwrap();
    // if !amount_reg.is_match(&amount) {
    //     return Err(format!("amount is {} not num", amount));
    // }
    let amount = match amount.parse::<u32>() {
        Ok(num) => num,
        Err(err) => {
            return Err(format!(
                "Maybe, the amount[{}] is too large.{:?}",
                amount, err
            ));
        }
    };
    Ok(amount)
}

// 将&[u8] byte转换为u32
fn get_u32_from_str_byte(amount: &[u8]) -> u32 {
    if amount.len() == 0 {
        return 0;
    }
    let amount = std::str::from_utf8(amount).unwrap();
    amount.parse::<u32>().unwrap()
}
// 根据公钥计算地址
// fn cal_address(pub_key: &str) -> String {
//     use data_encoding::HEXUPPER;
//     use ring::digest::{Context, SHA256};
//     let data = pub_key.as_bytes();

//     let mut context = Context::new(&SHA256);
//     context.update(data);
//     let digest = context.finish();

//     let hash = digest.as_ref();
//     let hex = HEXUPPER.encode(&hash[..160]);
//     hex
// }

// 根据公钥计算地址
fn calc_address(pub_key: &str) -> String {
    return pub_key.to_string();
}

