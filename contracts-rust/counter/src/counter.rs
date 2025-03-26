///
/// Copyright (C) BABEC. All rights reserved.
/// 
/// SPDX-License-Identifier: Apache-2.0
/// 

use contract_sdk_rust::sim_context;
use sim_context::*;

// 安装合约时会执行此方法，必须
#[no_mangle]
pub extern "C" fn init_contract() {
    // 安装时的业务逻辑，可为空
    sim_context::log("init_contract");
}

// 升级合约时会执行此方法，必须
#[no_mangle]
pub extern "C" fn upgrade() {
    // 升级时的业务逻辑，可为空
    sim_context::log("upgrade success");
}

#[no_mangle]
pub extern "C" fn increase() {
    sim_context::log("invoke increase");
    let ctx = &mut sim_context::get_sim_context();

    let mut key = ctx.arg_as_utf8_str("key");
    if key.len() == 0 {
        key = "counter1".to_string();
    }
    let r = ctx.get_state("count", &key);

    let count = r.unwrap();
    let counti32 = i32_from(count.as_slice());
    let count_new_i32 = counti32 + 1;
    let count_new = count_new_i32.to_le_bytes().to_vec();
    let is_ok = ctx.put_state("count", &key, count_new.as_slice());
    if is_ok == sim_context::SUCCESS_CODE {
        let s = format!("++ stone success count={}", count_new_i32.to_string());
        ctx.log(&s);
        let _ = ctx.ok(s.as_bytes());
    } else {
        let s = format!("++ stone fail old count={}", counti32.to_string());
        ctx.log(&s);
        let _ = ctx.error(&key);
    }
}

#[no_mangle]
pub extern "C" fn query() {
    sim_context::log("invoke query");
    let ctx = &mut sim_context::get_sim_context();

    let mut key = ctx.arg_as_utf8_str("key");
    if key.len() == 0 {
        key = "counter1".to_string();
    }

    let r = ctx.get_state("count", &key);
    let count = r.unwrap();
    let count: i32 = i32_from(count.as_slice());

    let log_str = format!("query count {}", count);
    let r = count.to_string();

    ctx.log(&log_str);
    ctx.ok(r.as_bytes());
}

pub fn i32_from(amount: &[u8]) -> i32 {
    if amount.len() == 4 {
        // receiver account money bytes into i32
        let b: i32 = i32::from_le_bytes([
            *amount.get(0).unwrap(),
            *amount.get(1).unwrap(),
            *amount.get(2).unwrap(),
            *amount.get(3).unwrap(),
        ]);
        return b;
    } else {
        return 0;
    }
}
