use parity_scale_codec;

fn main() {
    println!("Hello, world!");
}

#[cfg(test)]
mod tests {
    use super::parity_scale_codec::Encode;

    #[test]
    fn scale_encoded_u128() {
        let value: u128 = u64::MAX as u128;
        println!("{:?} - {:?}", value, value.encode());

        println!("Hello, world!");
        assert!(true);
    }

    #[test]
    fn scale_encoded_integers() {
        println!("{:?} - {:?}", u8::MAX, u8::MAX.encode());
        println!("{:?} - {:?}", i8::MAX, i8::MAX.encode());
        println!("{:?} - {:?}", -10 as i8, (-10 as i8).encode());

        println!("{:?} - {:?}", u16::MAX, u16::MAX.encode());
        println!("{:?} - {:?}", i16::MIN, i16::MIN.encode());
        println!("{:?} - {:?}", u32::MAX, u32::MAX.encode());
        println!("{:?} - {:?}", i32::MIN, i32::MIN.encode());
        println!("{:?} - {:?}", u64::MAX, u64::MAX.encode());
        println!("{:?} - {:?}", i64::MIN, i64::MIN.encode());
    }

    #[test]
    fn scale_encode_result() {
        let a: Result<u64, bool> = Ok(332290);
        println!("{:?}", a.encode());
        let b: Result<Result<u64, u8>, bool> = Ok(Err(10 as u8));
        println!("{:?}", b.encode());
    }

    #[test]
    fn scale_encode_option() {
        let a: Option<Option<bool>> = Some(None);
        println!("{:?}", a.encode());

        let a: Option<Result<String, bool>> = Some(Ok(String::from("eclesio")));
        println!("{:?}", a.encode());
    }
}
