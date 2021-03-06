package aesr

// Round32sle mixes the input values with aes tables and returns the result.
func Round32sle(x0, x1, x2, x3 uint32) (uint32, uint32, uint32, uint32) {
	y0 := (kAes0[x0&0xFF] ^
		kAes1[(x1>>8)&0xFF] ^
		kAes2[(x2>>16)&0xFF] ^
		kAes3[(x3>>24)&0xFF])
	y1 := (kAes0[x1&0xFF] ^
		kAes1[(x2>>8)&0xFF] ^
		kAes2[(x3>>16)&0xFF] ^
		kAes3[(x0>>24)&0xFF])
	y2 := (kAes0[x2&0xFF] ^
		kAes1[(x3>>8)&0xFF] ^
		kAes2[(x0>>16)&0xFF] ^
		kAes3[(x1>>24)&0xFF])
	y3 := (kAes0[x3&0xFF] ^
		kAes1[(x0>>8)&0xFF] ^
		kAes2[(x1>>16)&0xFF] ^
		kAes3[(x2>>24)&0xFF])
	return y0, y1, y2, y3
}

// Round32ble mixes the input values with aes tables and returns the result.
func Round32ble(x0, x1, x2, x3, k0, k1, k2, k3 uint32) (uint32, uint32, uint32, uint32) {
	y0 := (kAes0[x0&0xFF] ^
		kAes1[(x1>>8)&0xFF] ^
		kAes2[(x2>>16)&0xFF] ^
		kAes3[(x3>>24)&0xFF] ^ k0)
	y1 := (kAes0[x1&0xFF] ^
		kAes1[(x2>>8)&0xFF] ^
		kAes2[(x3>>16)&0xFF] ^
		kAes3[(x0>>24)&0xFF] ^ k1)
	y2 := (kAes0[x2&0xFF] ^
		kAes1[(x3>>8)&0xFF] ^
		kAes2[(x0>>16)&0xFF] ^
		kAes3[(x1>>24)&0xFF] ^ k2)
	y3 := (kAes0[x3&0xFF] ^
		kAes1[(x0>>8)&0xFF] ^
		kAes2[(x1>>16)&0xFF] ^
		kAes3[(x2>>24)&0xFF] ^ k3)
	return y0, y1, y2, y3
}

var kAes0 = [256]uint32{
	uint32(0xA56363C6), uint32(0x847C7CF8), uint32(0x997777EE), uint32(0x8D7B7BF6),
	uint32(0x0DF2F2FF), uint32(0xBD6B6BD6), uint32(0xB16F6FDE), uint32(0x54C5C591),
	uint32(0x50303060), uint32(0x03010102), uint32(0xA96767CE), uint32(0x7D2B2B56),
	uint32(0x19FEFEE7), uint32(0x62D7D7B5), uint32(0xE6ABAB4D), uint32(0x9A7676EC),
	uint32(0x45CACA8F), uint32(0x9D82821F), uint32(0x40C9C989), uint32(0x877D7DFA),
	uint32(0x15FAFAEF), uint32(0xEB5959B2), uint32(0xC947478E), uint32(0x0BF0F0FB),
	uint32(0xECADAD41), uint32(0x67D4D4B3), uint32(0xFDA2A25F), uint32(0xEAAFAF45),
	uint32(0xBF9C9C23), uint32(0xF7A4A453), uint32(0x967272E4), uint32(0x5BC0C09B),
	uint32(0xC2B7B775), uint32(0x1CFDFDE1), uint32(0xAE93933D), uint32(0x6A26264C),
	uint32(0x5A36366C), uint32(0x413F3F7E), uint32(0x02F7F7F5), uint32(0x4FCCCC83),
	uint32(0x5C343468), uint32(0xF4A5A551), uint32(0x34E5E5D1), uint32(0x08F1F1F9),
	uint32(0x937171E2), uint32(0x73D8D8AB), uint32(0x53313162), uint32(0x3F15152A),
	uint32(0x0C040408), uint32(0x52C7C795), uint32(0x65232346), uint32(0x5EC3C39D),
	uint32(0x28181830), uint32(0xA1969637), uint32(0x0F05050A), uint32(0xB59A9A2F),
	uint32(0x0907070E), uint32(0x36121224), uint32(0x9B80801B), uint32(0x3DE2E2DF),
	uint32(0x26EBEBCD), uint32(0x6927274E), uint32(0xCDB2B27F), uint32(0x9F7575EA),
	uint32(0x1B090912), uint32(0x9E83831D), uint32(0x742C2C58), uint32(0x2E1A1A34),
	uint32(0x2D1B1B36), uint32(0xB26E6EDC), uint32(0xEE5A5AB4), uint32(0xFBA0A05B),
	uint32(0xF65252A4), uint32(0x4D3B3B76), uint32(0x61D6D6B7), uint32(0xCEB3B37D),
	uint32(0x7B292952), uint32(0x3EE3E3DD), uint32(0x712F2F5E), uint32(0x97848413),
	uint32(0xF55353A6), uint32(0x68D1D1B9), uint32(0x00000000), uint32(0x2CEDEDC1),
	uint32(0x60202040), uint32(0x1FFCFCE3), uint32(0xC8B1B179), uint32(0xED5B5BB6),
	uint32(0xBE6A6AD4), uint32(0x46CBCB8D), uint32(0xD9BEBE67), uint32(0x4B393972),
	uint32(0xDE4A4A94), uint32(0xD44C4C98), uint32(0xE85858B0), uint32(0x4ACFCF85),
	uint32(0x6BD0D0BB), uint32(0x2AEFEFC5), uint32(0xE5AAAA4F), uint32(0x16FBFBED),
	uint32(0xC5434386), uint32(0xD74D4D9A), uint32(0x55333366), uint32(0x94858511),
	uint32(0xCF45458A), uint32(0x10F9F9E9), uint32(0x06020204), uint32(0x817F7FFE),
	uint32(0xF05050A0), uint32(0x443C3C78), uint32(0xBA9F9F25), uint32(0xE3A8A84B),
	uint32(0xF35151A2), uint32(0xFEA3A35D), uint32(0xC0404080), uint32(0x8A8F8F05),
	uint32(0xAD92923F), uint32(0xBC9D9D21), uint32(0x48383870), uint32(0x04F5F5F1),
	uint32(0xDFBCBC63), uint32(0xC1B6B677), uint32(0x75DADAAF), uint32(0x63212142),
	uint32(0x30101020), uint32(0x1AFFFFE5), uint32(0x0EF3F3FD), uint32(0x6DD2D2BF),
	uint32(0x4CCDCD81), uint32(0x140C0C18), uint32(0x35131326), uint32(0x2FECECC3),
	uint32(0xE15F5FBE), uint32(0xA2979735), uint32(0xCC444488), uint32(0x3917172E),
	uint32(0x57C4C493), uint32(0xF2A7A755), uint32(0x827E7EFC), uint32(0x473D3D7A),
	uint32(0xAC6464C8), uint32(0xE75D5DBA), uint32(0x2B191932), uint32(0x957373E6),
	uint32(0xA06060C0), uint32(0x98818119), uint32(0xD14F4F9E), uint32(0x7FDCDCA3),
	uint32(0x66222244), uint32(0x7E2A2A54), uint32(0xAB90903B), uint32(0x8388880B),
	uint32(0xCA46468C), uint32(0x29EEEEC7), uint32(0xD3B8B86B), uint32(0x3C141428),
	uint32(0x79DEDEA7), uint32(0xE25E5EBC), uint32(0x1D0B0B16), uint32(0x76DBDBAD),
	uint32(0x3BE0E0DB), uint32(0x56323264), uint32(0x4E3A3A74), uint32(0x1E0A0A14),
	uint32(0xDB494992), uint32(0x0A06060C), uint32(0x6C242448), uint32(0xE45C5CB8),
	uint32(0x5DC2C29F), uint32(0x6ED3D3BD), uint32(0xEFACAC43), uint32(0xA66262C4),
	uint32(0xA8919139), uint32(0xA4959531), uint32(0x37E4E4D3), uint32(0x8B7979F2),
	uint32(0x32E7E7D5), uint32(0x43C8C88B), uint32(0x5937376E), uint32(0xB76D6DDA),
	uint32(0x8C8D8D01), uint32(0x64D5D5B1), uint32(0xD24E4E9C), uint32(0xE0A9A949),
	uint32(0xB46C6CD8), uint32(0xFA5656AC), uint32(0x07F4F4F3), uint32(0x25EAEACF),
	uint32(0xAF6565CA), uint32(0x8E7A7AF4), uint32(0xE9AEAE47), uint32(0x18080810),
	uint32(0xD5BABA6F), uint32(0x887878F0), uint32(0x6F25254A), uint32(0x722E2E5C),
	uint32(0x241C1C38), uint32(0xF1A6A657), uint32(0xC7B4B473), uint32(0x51C6C697),
	uint32(0x23E8E8CB), uint32(0x7CDDDDA1), uint32(0x9C7474E8), uint32(0x211F1F3E),
	uint32(0xDD4B4B96), uint32(0xDCBDBD61), uint32(0x868B8B0D), uint32(0x858A8A0F),
	uint32(0x907070E0), uint32(0x423E3E7C), uint32(0xC4B5B571), uint32(0xAA6666CC),
	uint32(0xD8484890), uint32(0x05030306), uint32(0x01F6F6F7), uint32(0x120E0E1C),
	uint32(0xA36161C2), uint32(0x5F35356A), uint32(0xF95757AE), uint32(0xD0B9B969),
	uint32(0x91868617), uint32(0x58C1C199), uint32(0x271D1D3A), uint32(0xB99E9E27),
	uint32(0x38E1E1D9), uint32(0x13F8F8EB), uint32(0xB398982B), uint32(0x33111122),
	uint32(0xBB6969D2), uint32(0x70D9D9A9), uint32(0x898E8E07), uint32(0xA7949433),
	uint32(0xB69B9B2D), uint32(0x221E1E3C), uint32(0x92878715), uint32(0x20E9E9C9),
	uint32(0x49CECE87), uint32(0xFF5555AA), uint32(0x78282850), uint32(0x7ADFDFA5),
	uint32(0x8F8C8C03), uint32(0xF8A1A159), uint32(0x80898909), uint32(0x170D0D1A),
	uint32(0xDABFBF65), uint32(0x31E6E6D7), uint32(0xC6424284), uint32(0xB86868D0),
	uint32(0xC3414182), uint32(0xB0999929), uint32(0x772D2D5A), uint32(0x110F0F1E),
	uint32(0xCBB0B07B), uint32(0xFC5454A8), uint32(0xD6BBBB6D), uint32(0x3A16162C),
}

var kAes1 = [256]uint32{
	uint32(0x6363C6A5), uint32(0x7C7CF884), uint32(0x7777EE99), uint32(0x7B7BF68D),
	uint32(0xF2F2FF0D), uint32(0x6B6BD6BD), uint32(0x6F6FDEB1), uint32(0xC5C59154),
	uint32(0x30306050), uint32(0x01010203), uint32(0x6767CEA9), uint32(0x2B2B567D),
	uint32(0xFEFEE719), uint32(0xD7D7B562), uint32(0xABAB4DE6), uint32(0x7676EC9A),
	uint32(0xCACA8F45), uint32(0x82821F9D), uint32(0xC9C98940), uint32(0x7D7DFA87),
	uint32(0xFAFAEF15), uint32(0x5959B2EB), uint32(0x47478EC9), uint32(0xF0F0FB0B),
	uint32(0xADAD41EC), uint32(0xD4D4B367), uint32(0xA2A25FFD), uint32(0xAFAF45EA),
	uint32(0x9C9C23BF), uint32(0xA4A453F7), uint32(0x7272E496), uint32(0xC0C09B5B),
	uint32(0xB7B775C2), uint32(0xFDFDE11C), uint32(0x93933DAE), uint32(0x26264C6A),
	uint32(0x36366C5A), uint32(0x3F3F7E41), uint32(0xF7F7F502), uint32(0xCCCC834F),
	uint32(0x3434685C), uint32(0xA5A551F4), uint32(0xE5E5D134), uint32(0xF1F1F908),
	uint32(0x7171E293), uint32(0xD8D8AB73), uint32(0x31316253), uint32(0x15152A3F),
	uint32(0x0404080C), uint32(0xC7C79552), uint32(0x23234665), uint32(0xC3C39D5E),
	uint32(0x18183028), uint32(0x969637A1), uint32(0x05050A0F), uint32(0x9A9A2FB5),
	uint32(0x07070E09), uint32(0x12122436), uint32(0x80801B9B), uint32(0xE2E2DF3D),
	uint32(0xEBEBCD26), uint32(0x27274E69), uint32(0xB2B27FCD), uint32(0x7575EA9F),
	uint32(0x0909121B), uint32(0x83831D9E), uint32(0x2C2C5874), uint32(0x1A1A342E),
	uint32(0x1B1B362D), uint32(0x6E6EDCB2), uint32(0x5A5AB4EE), uint32(0xA0A05BFB),
	uint32(0x5252A4F6), uint32(0x3B3B764D), uint32(0xD6D6B761), uint32(0xB3B37DCE),
	uint32(0x2929527B), uint32(0xE3E3DD3E), uint32(0x2F2F5E71), uint32(0x84841397),
	uint32(0x5353A6F5), uint32(0xD1D1B968), uint32(0x00000000), uint32(0xEDEDC12C),
	uint32(0x20204060), uint32(0xFCFCE31F), uint32(0xB1B179C8), uint32(0x5B5BB6ED),
	uint32(0x6A6AD4BE), uint32(0xCBCB8D46), uint32(0xBEBE67D9), uint32(0x3939724B),
	uint32(0x4A4A94DE), uint32(0x4C4C98D4), uint32(0x5858B0E8), uint32(0xCFCF854A),
	uint32(0xD0D0BB6B), uint32(0xEFEFC52A), uint32(0xAAAA4FE5), uint32(0xFBFBED16),
	uint32(0x434386C5), uint32(0x4D4D9AD7), uint32(0x33336655), uint32(0x85851194),
	uint32(0x45458ACF), uint32(0xF9F9E910), uint32(0x02020406), uint32(0x7F7FFE81),
	uint32(0x5050A0F0), uint32(0x3C3C7844), uint32(0x9F9F25BA), uint32(0xA8A84BE3),
	uint32(0x5151A2F3), uint32(0xA3A35DFE), uint32(0x404080C0), uint32(0x8F8F058A),
	uint32(0x92923FAD), uint32(0x9D9D21BC), uint32(0x38387048), uint32(0xF5F5F104),
	uint32(0xBCBC63DF), uint32(0xB6B677C1), uint32(0xDADAAF75), uint32(0x21214263),
	uint32(0x10102030), uint32(0xFFFFE51A), uint32(0xF3F3FD0E), uint32(0xD2D2BF6D),
	uint32(0xCDCD814C), uint32(0x0C0C1814), uint32(0x13132635), uint32(0xECECC32F),
	uint32(0x5F5FBEE1), uint32(0x979735A2), uint32(0x444488CC), uint32(0x17172E39),
	uint32(0xC4C49357), uint32(0xA7A755F2), uint32(0x7E7EFC82), uint32(0x3D3D7A47),
	uint32(0x6464C8AC), uint32(0x5D5DBAE7), uint32(0x1919322B), uint32(0x7373E695),
	uint32(0x6060C0A0), uint32(0x81811998), uint32(0x4F4F9ED1), uint32(0xDCDCA37F),
	uint32(0x22224466), uint32(0x2A2A547E), uint32(0x90903BAB), uint32(0x88880B83),
	uint32(0x46468CCA), uint32(0xEEEEC729), uint32(0xB8B86BD3), uint32(0x1414283C),
	uint32(0xDEDEA779), uint32(0x5E5EBCE2), uint32(0x0B0B161D), uint32(0xDBDBAD76),
	uint32(0xE0E0DB3B), uint32(0x32326456), uint32(0x3A3A744E), uint32(0x0A0A141E),
	uint32(0x494992DB), uint32(0x06060C0A), uint32(0x2424486C), uint32(0x5C5CB8E4),
	uint32(0xC2C29F5D), uint32(0xD3D3BD6E), uint32(0xACAC43EF), uint32(0x6262C4A6),
	uint32(0x919139A8), uint32(0x959531A4), uint32(0xE4E4D337), uint32(0x7979F28B),
	uint32(0xE7E7D532), uint32(0xC8C88B43), uint32(0x37376E59), uint32(0x6D6DDAB7),
	uint32(0x8D8D018C), uint32(0xD5D5B164), uint32(0x4E4E9CD2), uint32(0xA9A949E0),
	uint32(0x6C6CD8B4), uint32(0x5656ACFA), uint32(0xF4F4F307), uint32(0xEAEACF25),
	uint32(0x6565CAAF), uint32(0x7A7AF48E), uint32(0xAEAE47E9), uint32(0x08081018),
	uint32(0xBABA6FD5), uint32(0x7878F088), uint32(0x25254A6F), uint32(0x2E2E5C72),
	uint32(0x1C1C3824), uint32(0xA6A657F1), uint32(0xB4B473C7), uint32(0xC6C69751),
	uint32(0xE8E8CB23), uint32(0xDDDDA17C), uint32(0x7474E89C), uint32(0x1F1F3E21),
	uint32(0x4B4B96DD), uint32(0xBDBD61DC), uint32(0x8B8B0D86), uint32(0x8A8A0F85),
	uint32(0x7070E090), uint32(0x3E3E7C42), uint32(0xB5B571C4), uint32(0x6666CCAA),
	uint32(0x484890D8), uint32(0x03030605), uint32(0xF6F6F701), uint32(0x0E0E1C12),
	uint32(0x6161C2A3), uint32(0x35356A5F), uint32(0x5757AEF9), uint32(0xB9B969D0),
	uint32(0x86861791), uint32(0xC1C19958), uint32(0x1D1D3A27), uint32(0x9E9E27B9),
	uint32(0xE1E1D938), uint32(0xF8F8EB13), uint32(0x98982BB3), uint32(0x11112233),
	uint32(0x6969D2BB), uint32(0xD9D9A970), uint32(0x8E8E0789), uint32(0x949433A7),
	uint32(0x9B9B2DB6), uint32(0x1E1E3C22), uint32(0x87871592), uint32(0xE9E9C920),
	uint32(0xCECE8749), uint32(0x5555AAFF), uint32(0x28285078), uint32(0xDFDFA57A),
	uint32(0x8C8C038F), uint32(0xA1A159F8), uint32(0x89890980), uint32(0x0D0D1A17),
	uint32(0xBFBF65DA), uint32(0xE6E6D731), uint32(0x424284C6), uint32(0x6868D0B8),
	uint32(0x414182C3), uint32(0x999929B0), uint32(0x2D2D5A77), uint32(0x0F0F1E11),
	uint32(0xB0B07BCB), uint32(0x5454A8FC), uint32(0xBBBB6DD6), uint32(0x16162C3A),
}

var kAes2 = [256]uint32{
	uint32(0x63C6A563), uint32(0x7CF8847C), uint32(0x77EE9977), uint32(0x7BF68D7B),
	uint32(0xF2FF0DF2), uint32(0x6BD6BD6B), uint32(0x6FDEB16F), uint32(0xC59154C5),
	uint32(0x30605030), uint32(0x01020301), uint32(0x67CEA967), uint32(0x2B567D2B),
	uint32(0xFEE719FE), uint32(0xD7B562D7), uint32(0xAB4DE6AB), uint32(0x76EC9A76),
	uint32(0xCA8F45CA), uint32(0x821F9D82), uint32(0xC98940C9), uint32(0x7DFA877D),
	uint32(0xFAEF15FA), uint32(0x59B2EB59), uint32(0x478EC947), uint32(0xF0FB0BF0),
	uint32(0xAD41ECAD), uint32(0xD4B367D4), uint32(0xA25FFDA2), uint32(0xAF45EAAF),
	uint32(0x9C23BF9C), uint32(0xA453F7A4), uint32(0x72E49672), uint32(0xC09B5BC0),
	uint32(0xB775C2B7), uint32(0xFDE11CFD), uint32(0x933DAE93), uint32(0x264C6A26),
	uint32(0x366C5A36), uint32(0x3F7E413F), uint32(0xF7F502F7), uint32(0xCC834FCC),
	uint32(0x34685C34), uint32(0xA551F4A5), uint32(0xE5D134E5), uint32(0xF1F908F1),
	uint32(0x71E29371), uint32(0xD8AB73D8), uint32(0x31625331), uint32(0x152A3F15),
	uint32(0x04080C04), uint32(0xC79552C7), uint32(0x23466523), uint32(0xC39D5EC3),
	uint32(0x18302818), uint32(0x9637A196), uint32(0x050A0F05), uint32(0x9A2FB59A),
	uint32(0x070E0907), uint32(0x12243612), uint32(0x801B9B80), uint32(0xE2DF3DE2),
	uint32(0xEBCD26EB), uint32(0x274E6927), uint32(0xB27FCDB2), uint32(0x75EA9F75),
	uint32(0x09121B09), uint32(0x831D9E83), uint32(0x2C58742C), uint32(0x1A342E1A),
	uint32(0x1B362D1B), uint32(0x6EDCB26E), uint32(0x5AB4EE5A), uint32(0xA05BFBA0),
	uint32(0x52A4F652), uint32(0x3B764D3B), uint32(0xD6B761D6), uint32(0xB37DCEB3),
	uint32(0x29527B29), uint32(0xE3DD3EE3), uint32(0x2F5E712F), uint32(0x84139784),
	uint32(0x53A6F553), uint32(0xD1B968D1), uint32(0x00000000), uint32(0xEDC12CED),
	uint32(0x20406020), uint32(0xFCE31FFC), uint32(0xB179C8B1), uint32(0x5BB6ED5B),
	uint32(0x6AD4BE6A), uint32(0xCB8D46CB), uint32(0xBE67D9BE), uint32(0x39724B39),
	uint32(0x4A94DE4A), uint32(0x4C98D44C), uint32(0x58B0E858), uint32(0xCF854ACF),
	uint32(0xD0BB6BD0), uint32(0xEFC52AEF), uint32(0xAA4FE5AA), uint32(0xFBED16FB),
	uint32(0x4386C543), uint32(0x4D9AD74D), uint32(0x33665533), uint32(0x85119485),
	uint32(0x458ACF45), uint32(0xF9E910F9), uint32(0x02040602), uint32(0x7FFE817F),
	uint32(0x50A0F050), uint32(0x3C78443C), uint32(0x9F25BA9F), uint32(0xA84BE3A8),
	uint32(0x51A2F351), uint32(0xA35DFEA3), uint32(0x4080C040), uint32(0x8F058A8F),
	uint32(0x923FAD92), uint32(0x9D21BC9D), uint32(0x38704838), uint32(0xF5F104F5),
	uint32(0xBC63DFBC), uint32(0xB677C1B6), uint32(0xDAAF75DA), uint32(0x21426321),
	uint32(0x10203010), uint32(0xFFE51AFF), uint32(0xF3FD0EF3), uint32(0xD2BF6DD2),
	uint32(0xCD814CCD), uint32(0x0C18140C), uint32(0x13263513), uint32(0xECC32FEC),
	uint32(0x5FBEE15F), uint32(0x9735A297), uint32(0x4488CC44), uint32(0x172E3917),
	uint32(0xC49357C4), uint32(0xA755F2A7), uint32(0x7EFC827E), uint32(0x3D7A473D),
	uint32(0x64C8AC64), uint32(0x5DBAE75D), uint32(0x19322B19), uint32(0x73E69573),
	uint32(0x60C0A060), uint32(0x81199881), uint32(0x4F9ED14F), uint32(0xDCA37FDC),
	uint32(0x22446622), uint32(0x2A547E2A), uint32(0x903BAB90), uint32(0x880B8388),
	uint32(0x468CCA46), uint32(0xEEC729EE), uint32(0xB86BD3B8), uint32(0x14283C14),
	uint32(0xDEA779DE), uint32(0x5EBCE25E), uint32(0x0B161D0B), uint32(0xDBAD76DB),
	uint32(0xE0DB3BE0), uint32(0x32645632), uint32(0x3A744E3A), uint32(0x0A141E0A),
	uint32(0x4992DB49), uint32(0x060C0A06), uint32(0x24486C24), uint32(0x5CB8E45C),
	uint32(0xC29F5DC2), uint32(0xD3BD6ED3), uint32(0xAC43EFAC), uint32(0x62C4A662),
	uint32(0x9139A891), uint32(0x9531A495), uint32(0xE4D337E4), uint32(0x79F28B79),
	uint32(0xE7D532E7), uint32(0xC88B43C8), uint32(0x376E5937), uint32(0x6DDAB76D),
	uint32(0x8D018C8D), uint32(0xD5B164D5), uint32(0x4E9CD24E), uint32(0xA949E0A9),
	uint32(0x6CD8B46C), uint32(0x56ACFA56), uint32(0xF4F307F4), uint32(0xEACF25EA),
	uint32(0x65CAAF65), uint32(0x7AF48E7A), uint32(0xAE47E9AE), uint32(0x08101808),
	uint32(0xBA6FD5BA), uint32(0x78F08878), uint32(0x254A6F25), uint32(0x2E5C722E),
	uint32(0x1C38241C), uint32(0xA657F1A6), uint32(0xB473C7B4), uint32(0xC69751C6),
	uint32(0xE8CB23E8), uint32(0xDDA17CDD), uint32(0x74E89C74), uint32(0x1F3E211F),
	uint32(0x4B96DD4B), uint32(0xBD61DCBD), uint32(0x8B0D868B), uint32(0x8A0F858A),
	uint32(0x70E09070), uint32(0x3E7C423E), uint32(0xB571C4B5), uint32(0x66CCAA66),
	uint32(0x4890D848), uint32(0x03060503), uint32(0xF6F701F6), uint32(0x0E1C120E),
	uint32(0x61C2A361), uint32(0x356A5F35), uint32(0x57AEF957), uint32(0xB969D0B9),
	uint32(0x86179186), uint32(0xC19958C1), uint32(0x1D3A271D), uint32(0x9E27B99E),
	uint32(0xE1D938E1), uint32(0xF8EB13F8), uint32(0x982BB398), uint32(0x11223311),
	uint32(0x69D2BB69), uint32(0xD9A970D9), uint32(0x8E07898E), uint32(0x9433A794),
	uint32(0x9B2DB69B), uint32(0x1E3C221E), uint32(0x87159287), uint32(0xE9C920E9),
	uint32(0xCE8749CE), uint32(0x55AAFF55), uint32(0x28507828), uint32(0xDFA57ADF),
	uint32(0x8C038F8C), uint32(0xA159F8A1), uint32(0x89098089), uint32(0x0D1A170D),
	uint32(0xBF65DABF), uint32(0xE6D731E6), uint32(0x4284C642), uint32(0x68D0B868),
	uint32(0x4182C341), uint32(0x9929B099), uint32(0x2D5A772D), uint32(0x0F1E110F),
	uint32(0xB07BCBB0), uint32(0x54A8FC54), uint32(0xBB6DD6BB), uint32(0x162C3A16),
}

var kAes3 = [256]uint32{
	uint32(0xC6A56363), uint32(0xF8847C7C), uint32(0xEE997777), uint32(0xF68D7B7B),
	uint32(0xFF0DF2F2), uint32(0xD6BD6B6B), uint32(0xDEB16F6F), uint32(0x9154C5C5),
	uint32(0x60503030), uint32(0x02030101), uint32(0xCEA96767), uint32(0x567D2B2B),
	uint32(0xE719FEFE), uint32(0xB562D7D7), uint32(0x4DE6ABAB), uint32(0xEC9A7676),
	uint32(0x8F45CACA), uint32(0x1F9D8282), uint32(0x8940C9C9), uint32(0xFA877D7D),
	uint32(0xEF15FAFA), uint32(0xB2EB5959), uint32(0x8EC94747), uint32(0xFB0BF0F0),
	uint32(0x41ECADAD), uint32(0xB367D4D4), uint32(0x5FFDA2A2), uint32(0x45EAAFAF),
	uint32(0x23BF9C9C), uint32(0x53F7A4A4), uint32(0xE4967272), uint32(0x9B5BC0C0),
	uint32(0x75C2B7B7), uint32(0xE11CFDFD), uint32(0x3DAE9393), uint32(0x4C6A2626),
	uint32(0x6C5A3636), uint32(0x7E413F3F), uint32(0xF502F7F7), uint32(0x834FCCCC),
	uint32(0x685C3434), uint32(0x51F4A5A5), uint32(0xD134E5E5), uint32(0xF908F1F1),
	uint32(0xE2937171), uint32(0xAB73D8D8), uint32(0x62533131), uint32(0x2A3F1515),
	uint32(0x080C0404), uint32(0x9552C7C7), uint32(0x46652323), uint32(0x9D5EC3C3),
	uint32(0x30281818), uint32(0x37A19696), uint32(0x0A0F0505), uint32(0x2FB59A9A),
	uint32(0x0E090707), uint32(0x24361212), uint32(0x1B9B8080), uint32(0xDF3DE2E2),
	uint32(0xCD26EBEB), uint32(0x4E692727), uint32(0x7FCDB2B2), uint32(0xEA9F7575),
	uint32(0x121B0909), uint32(0x1D9E8383), uint32(0x58742C2C), uint32(0x342E1A1A),
	uint32(0x362D1B1B), uint32(0xDCB26E6E), uint32(0xB4EE5A5A), uint32(0x5BFBA0A0),
	uint32(0xA4F65252), uint32(0x764D3B3B), uint32(0xB761D6D6), uint32(0x7DCEB3B3),
	uint32(0x527B2929), uint32(0xDD3EE3E3), uint32(0x5E712F2F), uint32(0x13978484),
	uint32(0xA6F55353), uint32(0xB968D1D1), uint32(0x00000000), uint32(0xC12CEDED),
	uint32(0x40602020), uint32(0xE31FFCFC), uint32(0x79C8B1B1), uint32(0xB6ED5B5B),
	uint32(0xD4BE6A6A), uint32(0x8D46CBCB), uint32(0x67D9BEBE), uint32(0x724B3939),
	uint32(0x94DE4A4A), uint32(0x98D44C4C), uint32(0xB0E85858), uint32(0x854ACFCF),
	uint32(0xBB6BD0D0), uint32(0xC52AEFEF), uint32(0x4FE5AAAA), uint32(0xED16FBFB),
	uint32(0x86C54343), uint32(0x9AD74D4D), uint32(0x66553333), uint32(0x11948585),
	uint32(0x8ACF4545), uint32(0xE910F9F9), uint32(0x04060202), uint32(0xFE817F7F),
	uint32(0xA0F05050), uint32(0x78443C3C), uint32(0x25BA9F9F), uint32(0x4BE3A8A8),
	uint32(0xA2F35151), uint32(0x5DFEA3A3), uint32(0x80C04040), uint32(0x058A8F8F),
	uint32(0x3FAD9292), uint32(0x21BC9D9D), uint32(0x70483838), uint32(0xF104F5F5),
	uint32(0x63DFBCBC), uint32(0x77C1B6B6), uint32(0xAF75DADA), uint32(0x42632121),
	uint32(0x20301010), uint32(0xE51AFFFF), uint32(0xFD0EF3F3), uint32(0xBF6DD2D2),
	uint32(0x814CCDCD), uint32(0x18140C0C), uint32(0x26351313), uint32(0xC32FECEC),
	uint32(0xBEE15F5F), uint32(0x35A29797), uint32(0x88CC4444), uint32(0x2E391717),
	uint32(0x9357C4C4), uint32(0x55F2A7A7), uint32(0xFC827E7E), uint32(0x7A473D3D),
	uint32(0xC8AC6464), uint32(0xBAE75D5D), uint32(0x322B1919), uint32(0xE6957373),
	uint32(0xC0A06060), uint32(0x19988181), uint32(0x9ED14F4F), uint32(0xA37FDCDC),
	uint32(0x44662222), uint32(0x547E2A2A), uint32(0x3BAB9090), uint32(0x0B838888),
	uint32(0x8CCA4646), uint32(0xC729EEEE), uint32(0x6BD3B8B8), uint32(0x283C1414),
	uint32(0xA779DEDE), uint32(0xBCE25E5E), uint32(0x161D0B0B), uint32(0xAD76DBDB),
	uint32(0xDB3BE0E0), uint32(0x64563232), uint32(0x744E3A3A), uint32(0x141E0A0A),
	uint32(0x92DB4949), uint32(0x0C0A0606), uint32(0x486C2424), uint32(0xB8E45C5C),
	uint32(0x9F5DC2C2), uint32(0xBD6ED3D3), uint32(0x43EFACAC), uint32(0xC4A66262),
	uint32(0x39A89191), uint32(0x31A49595), uint32(0xD337E4E4), uint32(0xF28B7979),
	uint32(0xD532E7E7), uint32(0x8B43C8C8), uint32(0x6E593737), uint32(0xDAB76D6D),
	uint32(0x018C8D8D), uint32(0xB164D5D5), uint32(0x9CD24E4E), uint32(0x49E0A9A9),
	uint32(0xD8B46C6C), uint32(0xACFA5656), uint32(0xF307F4F4), uint32(0xCF25EAEA),
	uint32(0xCAAF6565), uint32(0xF48E7A7A), uint32(0x47E9AEAE), uint32(0x10180808),
	uint32(0x6FD5BABA), uint32(0xF0887878), uint32(0x4A6F2525), uint32(0x5C722E2E),
	uint32(0x38241C1C), uint32(0x57F1A6A6), uint32(0x73C7B4B4), uint32(0x9751C6C6),
	uint32(0xCB23E8E8), uint32(0xA17CDDDD), uint32(0xE89C7474), uint32(0x3E211F1F),
	uint32(0x96DD4B4B), uint32(0x61DCBDBD), uint32(0x0D868B8B), uint32(0x0F858A8A),
	uint32(0xE0907070), uint32(0x7C423E3E), uint32(0x71C4B5B5), uint32(0xCCAA6666),
	uint32(0x90D84848), uint32(0x06050303), uint32(0xF701F6F6), uint32(0x1C120E0E),
	uint32(0xC2A36161), uint32(0x6A5F3535), uint32(0xAEF95757), uint32(0x69D0B9B9),
	uint32(0x17918686), uint32(0x9958C1C1), uint32(0x3A271D1D), uint32(0x27B99E9E),
	uint32(0xD938E1E1), uint32(0xEB13F8F8), uint32(0x2BB39898), uint32(0x22331111),
	uint32(0xD2BB6969), uint32(0xA970D9D9), uint32(0x07898E8E), uint32(0x33A79494),
	uint32(0x2DB69B9B), uint32(0x3C221E1E), uint32(0x15928787), uint32(0xC920E9E9),
	uint32(0x8749CECE), uint32(0xAAFF5555), uint32(0x50782828), uint32(0xA57ADFDF),
	uint32(0x038F8C8C), uint32(0x59F8A1A1), uint32(0x09808989), uint32(0x1A170D0D),
	uint32(0x65DABFBF), uint32(0xD731E6E6), uint32(0x84C64242), uint32(0xD0B86868),
	uint32(0x82C34141), uint32(0x29B09999), uint32(0x5A772D2D), uint32(0x1E110F0F),
	uint32(0x7BCBB0B0), uint32(0xA8FC5454), uint32(0x6DD6BBBB), uint32(0x2C3A1616),
}
