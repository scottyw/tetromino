RangeTests:
	dw .all_bits_clear, AllBitsClearTest
	dw .all_bits_set, AllBitsSetTest
	dw .valid_bits, ValidBitsTest
	dw .invalid_value_tick | $8000, InvalidValueTickTest | $8000
	dw .invalid_rollovers | $8000, InvalidRolloversTest | $8000
	dw .high_minutes, HighMinutesTest
	dw .high_hours, HighHoursTest
	dw -1

.all_bits_clear
	db "All bits clear@"
.all_bits_set
	db "All bits set@"
.valid_bits
	db "Valid bits@"
.invalid_value_tick
	db "Invalid value tick@"
.invalid_rollovers
	db "Invalid rollovers@"
.high_minutes
	db "High minutes@"
.high_hours
	db "High hours@"

AllBitsClearTest:
	write_RTC_register RTCDH, 0
	; wait for a tick to prevent unrelated bugs from affecting this test
	call WaitNextRTCTick
	xor a
	ld b, a
	ld c, a
	ld d, a
	ld e, a
	call WriteRTC
	rst WaitVBlank
	call ReadRTC
	or b
	or c
	or d
	or e
	jr AllBitsSetTest.done

AllBitsSetTest:
	write_RTC_register RTCDH, $40
	ld a, $c1
	lb bc, $ff, $1f
	lb de, $3f, $3f
	call WriteRTC
	rst WaitVBlank
	call ReadRTC
	cp $c1
	jr nz, .done
	ld a, d
	cp e
	jr nz, .done
	cp $3f
	jr nz, .done
	inc b
	jr nz, .done
	ld a, c
	cp $1f
.done
	rst CarryIfNonZero
	jp PassFailResult

ValidBitsTest:
	; turn it off just in case
	write_RTC_register RTCDH, $40
	ld hl, rRAMB ;also initializes l = 0 for error tracking
	lb bc, $a0, $1f
	ld [hl], RTCH
	call .test
	rl l
	ld [hl], RTCM
	ld c, $3f
	call .test
	rl l
	ld [hl], RTCS
	call .test
	rl l
	; ensure the second counter is 0 so we don't accidentally get a rollover when turning the RTC on
	xor a
	ld [bc], a
	ld [hl], RTCDH
	ld c, $c1
	call .test
	scf
	and $10
	or l
	jp FailedRegistersResult

.test
	; in: c: mask
	call Random
	inc a
	jr z, .test
	dec a
	jr z, .test
	ld e, a
	cpl
	ld [bc], a
	and c
	ld d, a
	latch_RTC
	ld a, [bc]
	cp d
	scf
	ret nz
	ld a, e
	ld [bc], a
	and c
	ld e, a
	latch_RTC
	ld a, [bc]
	cp e
	rst CarryIfNonZero
	ret

InvalidValueTickTest:
	call Random
	or $fc
	inc a
	jr z, InvalidValueTickTest
	add a, 63 ;results in a value between 60 and 62
	ld e, a
	call Random
	and 3
	add a, 60
	ld d, a
	call Random
	and 7
	add a, 24
	ld c, a
	call Random
	ld b, a
	call Random
	and 1
	call WriteRTC
	inc e
	call WaitCompareRTC
	rst CarryIfNonZero
	jp PassFailResult

InvalidRolloversTest:
	ld l, 0
	call .random_minutes
	ld e, 63
	call WriteRTC
	ld e, l ;must be 0 here
	call .test
	call .random_hours
	lb de, 63, 59
	call WriteRTC
	ld de, 0
	call .test
	call .random_days
	lb de, 59, 59
	ld c, 31
	call WriteRTC
	ld de, 0
	ld c, d
	call .test
	ld a, l
	jp FailedRegistersResult

.random_minutes
	call Random
	and 63
	cp 59
	jr nc, .random_minutes
	ld d, a
.random_hours
	call Random
	and 31
	cp 23
	jr nc, .random_hours
	ld c, a
.random_days
	call Random
	ld b, a
	call Random
	and 1
	ret z
	inc b
	jr z, .random_days
	dec b
	ret

.test
	push hl
	call WaitCompareRTC
	rst CarryIfNonZero
	pop hl
	rl l
	ret

HighMinutesTest:
	call Random
	or $fc
	inc a
	jr z, HighMinutesTest
	add a, 63 ;results in a value between 60 and 62
	ld d, a
.hours
	call Random
	and $1f
	cp 23
	jr nc, .hours
	ld c, a
	call Random
	ld b, a
	call Random
	and 1
	ld e, 59
	call WriteRTC
	inc d
	ld e, 0
	jr HighHoursTest.done

HighHoursTest:
	call Random
	or $f8
	inc a
	jr z, HighHoursTest
	add a, 31 ;results in a value between 24 and 30
	ld c, a
.days
	call Random
	inc a
	jr z, .days
	dec a
	ld b, a
	call Random
	and 1
	lb de, 59, 59
	call WriteRTC
	ld de, 0
	inc c
.done
	call WaitCompareRTC
	rst CarryIfNonZero
	jp PassFailResult
