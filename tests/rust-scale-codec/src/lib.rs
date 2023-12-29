use parity_scale_codec;

fn main() {
    println!("Hello, world!");
}

#[cfg(test)]
mod tests {
    use parity_scale_codec::Decode;

    use super::parity_scale_codec::Encode;
    use super::parity_scale_codec::Encode as DervieEncode;

    #[derive(DervieEncode, Decode, Debug)]
    enum Nested {
        SimpleN(u32),
    }

    #[derive(DervieEncode, Decode, Debug)]
    enum Error {
        FailureX,
    }

    #[derive(DervieEncode, Decode, Debug)]
    enum EnumA {
        A(bool),
    }

    #[derive(DervieEncode, Decode, Debug)]
    enum ToTest {
        Single,
        Int(u64),
        Bool(bool),
        A(Option<bool>),
        B(Result<u64, u64>),
        G((u64, bool)),
        H(Option<(u64, bool)>),
        J(Result<(u64, bool), bool>),
        K((Option<bool>, Result<bool, bool>)),
        L(Result<Option<(u64, bool)>, u64>),
        M(Option<Nested>),
        N(Result<Nested, bool>),
        O(Result<bool, Nested>),
        P(Result<Nested, Error>),
        Q((Nested, u64, Error)),
    }

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
        let a: Result<u64, bool> = Ok(32290);
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

    #[test]
    fn scale_encode_tuple() {
        let a: (Option<u64>, Result<bool, u64>) = (Some(79), Ok(true));
        println!("{:?}", a.encode());

        let a: (Option<u64>, Result<bool, u64>) = (None, Err(44));
        println!("{:?}", a.encode());
    }

    #[test]
    fn scale_encode_enum() {
        // let a: ToTest = ToTest::Single;
        // println!("{:?}", a.encode());

        // let a: ToTest = ToTest::Int(32);
        // println!("{:?}", a.encode());

        // let a: ToTest = ToTest::Bool(true);
        // println!("{:?}", a.encode());

        // let a: ToTest = ToTest::A(Some(true));
        // println!("{:?}", a.encode());

        // let a: ToTest = ToTest::A(None);
        // println!("{:?}", a.encode());

        // let a: ToTest = ToTest::B(Ok(108));
        // println!("{:?}", a.encode());

        // let a: ToTest = ToTest::B(Err(90));
        // println!("{:?}", a.encode());

        // let a: ToTest = ToTest::G((60, false));
        // println!("{:?}", a.encode());

        // let a: ToTest = ToTest::H(Some((60, false)));
        // println!("{:?}", a.encode());

        // let a: ToTest = ToTest::J(Ok((60, false)));
        // println!("{:?}", a.encode());

        // let a: ToTest = ToTest::K((Some(true), Ok(false)));
        // println!("{:?}", a.encode());

        // let a: ToTest = ToTest::M(Some(Nested::SimpleN(10)));
        // println!("{:?}", a.encode());

        // let a: ToTest = ToTest::N(Ok(Nested::SimpleN(78)));
        // println!("{:?}", a.encode());

        // let a: ToTest = ToTest::N(Err(true));
        // println!("{:?}", a.encode());

        // let a: ToTest = ToTest::O(Ok(true));
        // println!("{:?}", a.encode());

        // let a: ToTest = ToTest::O(Err(Nested::SimpleN(76)));
        // println!("{:?}", a.encode());

        // let a: ToTest = ToTest::P(Ok(Nested::SimpleN(u32::MAX)));
        // println!("{:?}", a.encode());

        // let a: ToTest = ToTest::P(Err(Error::FailureX));
        // println!("{:?}", a.encode());

        // let a: ToTest = ToTest::Q((Nested::SimpleN(77), 89, Error::FailureX));
        // println!("{:?}", a.encode());

        let binding = vec![1,32,0,0,0,0,0,0,0];
        let mut b: &[u8] = binding.as_slice();
        let res = ToTest::decode(&mut b).unwrap();

        println!("{:?}", res);
    }
}
