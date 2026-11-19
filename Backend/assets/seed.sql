-- ==============================================================================
-- Fuel Management System (وين الوقود؟) - Full Database Schema & Seed Data
-- ==============================================================================

-- 1. Create Tables
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(150) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    phone VARCHAR(30) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'user',
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    loyalty_points BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS stations (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(180) NOT NULL,
    location VARCHAR(255) NOT NULL,
    bank_account_number VARCHAR(100) NOT NULL,
    bank_account_name VARCHAR(180) NOT NULL,
    bank_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'pending',
    reviewed_by BIGINT,
    reviewed_at TIMESTAMPTZ,
    rejection_reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS branches (
    id BIGSERIAL PRIMARY KEY,
    station_id BIGINT NOT NULL REFERENCES stations(id) ON DELETE CASCADE,
    city VARCHAR(120) NOT NULL,
    address VARCHAR(255) NOT NULL,
    is_main BOOLEAN NOT NULL DEFAULT false,
    fuel_status VARCHAR(20) NOT NULL DEFAULT 'full',
    congestion_percentage INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS wallets (
    id BIGSERIAL PRIMARY KEY,
    owner_type VARCHAR(20) NOT NULL,
    owner_id BIGINT NOT NULL,
    total_balance BIGINT NOT NULL DEFAULT 0,
    allocated_balance BIGINT NOT NULL DEFAULT 0,
    available_balance BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS transactions (
    id BIGSERIAL PRIMARY KEY,
    reference VARCHAR(80) NOT NULL UNIQUE,
    actor_key VARCHAR(60),
    idempotency_key VARCHAR(120),
    request_hash VARCHAR(64),
    type VARCHAR(40) NOT NULL,
    amount BIGINT NOT NULL,
    from_type VARCHAR(20),
    from_id BIGINT,
    from_name VARCHAR(180) DEFAULT '',
    to_type VARCHAR(20),
    to_id BIGINT,
    to_name VARCHAR(180) DEFAULT '',
    branch_id BIGINT,
    branch_name VARCHAR(180),
    loyalty_points_used BIGINT NOT NULL DEFAULT 0,
    loyalty_discount BIGINT NOT NULL DEFAULT 0,
    status VARCHAR(30) NOT NULL DEFAULT 'completed',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS cash_withdrawals (
    id BIGSERIAL PRIMARY KEY,
    station_id BIGINT NOT NULL REFERENCES stations(id) ON DELETE CASCADE,
    idempotency_key VARCHAR(120),
    request_hash VARCHAR(64),
    amount BIGINT NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'pending',
    external_reference VARCHAR(120),
    reviewed_by BIGINT,
    reviewed_at TIMESTAMPTZ,
    rejection_reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    token_hash TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 2. Constraints & Indexes
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_wallets_non_negative') THEN
        ALTER TABLE wallets
        ADD CONSTRAINT chk_wallets_non_negative
        CHECK (available_balance >= 0 AND allocated_balance >= 0 AND total_balance >= 0);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'chk_users_loyalty_non_negative') THEN
        ALTER TABLE users
        ADD CONSTRAINT chk_users_loyalty_non_negative
        CHECK (loyalty_points >= 0);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_wallets_owner ON wallets(owner_type, owner_id);
CREATE INDEX IF NOT EXISTS idx_branches_city ON branches(city);
CREATE INDEX IF NOT EXISTS idx_branches_station ON branches(station_id);
CREATE INDEX IF NOT EXISTS idx_cw_station ON cash_withdrawals(station_id);
CREATE INDEX IF NOT EXISTS idx_tx_to ON transactions(to_id);
CREATE INDEX IF NOT EXISTS idx_tx_from ON transactions(from_id);

-- 3. Truncate & Clean Existing Data
TRUNCATE TABLE refresh_tokens, cash_withdrawals, transactions, wallets, branches, stations, users RESTART IDENTITY CASCADE;

-- 4. Seed Admins & Regular Users (Password for all: Password123!)
-- Hash: $2b$10$7jptClX.y4Wn58csfxYVJ.PqvbIkHz10bU2CuA.fUmn1Xdfq0KCbK
INSERT INTO users (id, name, email, password_hash, phone, role, status, loyalty_points, created_at) VALUES
(1, 'مدير النظام', 'admin@fuelms.com', '$2b$10$7jptClX.y4Wn58csfxYVJ.PqvbIkHz10bU2CuA.fUmn1Xdfq0KCbK', '+249912300001', 'admin', 'active', 0, NOW() - INTERVAL '30 days'),
(2, 'مشرف العمليات', 'supervisor@fuelms.com', '$2b$10$7jptClX.y4Wn58csfxYVJ.PqvbIkHz10bU2CuA.fUmn1Xdfq0KCbK', '+249912300002', 'admin', 'active', 0, NOW() - INTERVAL '28 days'),
(3, 'أحمد خليل', 'ahmed.khalil@example.com', '$2b$10$7jptClX.y4Wn58csfxYVJ.PqvbIkHz10bU2CuA.fUmn1Xdfq0KCbK', '+249912345601', 'user', 'active', 15, NOW() - INTERVAL '20 days'),
(4, 'محمد نور', 'mohamed.nour@example.com', '$2b$10$7jptClX.y4Wn58csfxYVJ.PqvbIkHz10bU2CuA.fUmn1Xdfq0KCbK', '+249912345602', 'user', 'active', 15, NOW() - INTERVAL '19 days'),
(5, 'فاطمة آدم', 'fatima.adam@example.com', '$2b$10$7jptClX.y4Wn58csfxYVJ.PqvbIkHz10bU2CuA.fUmn1Xdfq0KCbK', '+249912345603', 'user', 'active', 15, NOW() - INTERVAL '18 days'),
(6, 'عمر حسن', 'omer.hassan@example.com', '$2b$10$7jptClX.y4Wn58csfxYVJ.PqvbIkHz10bU2CuA.fUmn1Xdfq0KCbK', '+249912345604', 'user', 'active', 15, NOW() - INTERVAL '17 days'),
(7, 'سارة إبراهيم', 'sara.ibrahim@example.com', '$2b$10$7jptClX.y4Wn58csfxYVJ.PqvbIkHz10bU2CuA.fUmn1Xdfq0KCbK', '+249912345605', 'user', 'active', 15, NOW() - INTERVAL '16 days'),
(8, 'علي عثمان', 'ali.osman@example.com', '$2b$10$7jptClX.y4Wn58csfxYVJ.PqvbIkHz10bU2CuA.fUmn1Xdfq0KCbK', '+249912345606', 'user', 'active', 15, NOW() - INTERVAL '15 days'),
(9, 'منى صديق', 'mona.siddig@example.com', '$2b$10$7jptClX.y4Wn58csfxYVJ.PqvbIkHz10bU2CuA.fUmn1Xdfq0KCbK', '+249912345607', 'user', 'active', 15, NOW() - INTERVAL '14 days'),
(10, 'خالد إدريس', 'khalid.idris@example.com', '$2b$10$7jptClX.y4Wn58csfxYVJ.PqvbIkHz10bU2CuA.fUmn1Xdfq0KCbK', '+249912345608', 'user', 'active', 15, NOW() - INTERVAL '13 days'),
(11, 'زينب طارق', 'zainab.tariq@example.com', '$2b$10$7jptClX.y4Wn58csfxYVJ.PqvbIkHz10bU2CuA.fUmn1Xdfq0KCbK', '+249912345609', 'user', 'active', 15, NOW() - INTERVAL '12 days'),
(12, 'مصطفى بكري', 'mustafa.bakri@example.com', '$2b$10$7jptClX.y4Wn58csfxYVJ.PqvbIkHz10bU2CuA.fUmn1Xdfq0KCbK', '+249912345610', 'user', 'suspended', 15, NOW() - INTERVAL '11 days');

ALTER SEQUENCE users_id_seq RESTART WITH 13;

-- 5. Seed Stations (8 stations)
INSERT INTO stations (id, name, location, bank_account_number, bank_account_name, bank_name, email, password_hash, status, reviewed_by, reviewed_at, rejection_reason, created_at) VALUES
(1, 'شركة النيل للبترول', 'الخرطوم', '100293847501', 'شركة النيل للبترول', 'بنك الخرطوم', 'al-nile@petroleum.sd', '$2b$10$7jptClX.y4Wn58csfxYVJ.PqvbIkHz10bU2CuA.fUmn1Xdfq0KCbK', 'active', 1, NOW() - INTERVAL '25 days', NULL, NOW() - INTERVAL '26 days'),
(2, 'شركة بشائر للطاقة', 'الخرطوم', '200481928374', 'شركة بشائر للطاقة', 'بنك أم درمان الوطني', 'contact@bashair-energy.sd', '$2b$10$7jptClX.y4Wn58csfxYVJ.PqvbIkHz10bU2CuA.fUmn1Xdfq0KCbK', 'active', 1, NOW() - INTERVAL '25 days', NULL, NOW() - INTERVAL '25 days'),
(3, 'بتروتشاينا السودان', 'بحري', '300192837465', 'بتروتشاينا السودان', 'بنك فيصل الإسلامي', 'info@petrochina.sd', '$2b$10$7jptClX.y4Wn58csfxYVJ.PqvbIkHz10bU2CuA.fUmn1Xdfq0KCbK', 'active', 2, NOW() - INTERVAL '22 days', NULL, NOW() - INTERVAL '24 days'),
(4, 'النهضة لخدمات الوقود', 'أم درمان', '100847291039', 'النهضة لخدمات الوقود', 'بنك الخرطوم', 'support@alnahda-fuel.sd', '$2b$10$7jptClX.y4Wn58csfxYVJ.PqvbIkHz10bU2CuA.fUmn1Xdfq0KCbK', 'active', 1, NOW() - INTERVAL '22 days', NULL, NOW() - INTERVAL '23 days'),
(5, 'بترول البحر الأحمر', 'بورتسودان', '300994827163', 'بترول البحر الأحمر', 'بنك فيصل الإسلامي', 'operations@redsea-petro.sd', '$2b$10$7jptClX.y4Wn58csfxYVJ.PqvbIkHz10bU2CuA.fUmn1Xdfq0KCbK', 'active', 2, NOW() - INTERVAL '22 days', NULL, NOW() - INTERVAL '21 days'),
(6, 'الصحراء للنفط والغاز', 'الأبيض', '200771829304', 'شركة الصحراء للنفط والغاز المحدودة', 'بنك أم درمان الوطني', 'admin@sahara-oil.sd', '$2b$10$7jptClX.y4Wn58csfxYVJ.PqvbIkHz10bU2CuA.fUmn1Xdfq0KCbK', 'pending', NULL, NULL, NULL, NOW() - INTERVAL '2 days'),
(7, 'شركة النيل الأزرق للوقود', 'ود مدني', '100662718293', 'شركة النيل الأزرق للوقود المحدودة', 'بنك الخرطوم', 'info@bluenile-fuel.sd', '$2b$10$7jptClX.y4Wn58csfxYVJ.PqvbIkHz10bU2CuA.fUmn1Xdfq0KCbK', 'rejected', 1, NOW() - INTERVAL '5 days', 'انتهاء صلاحية السجل التجاري وعدم إمكانية التحقق من الضمان البنكي', NOW() - INTERVAL '7 days'),
(8, 'شركة تمازج للطاقة', 'كسلا', '300882910293', 'شركة تمازج للطاقة', 'بنك فيصل الإسلامي', 'contact@tamazuj-energy.sd', '$2b$10$7jptClX.y4Wn58csfxYVJ.PqvbIkHz10bU2CuA.fUmn1Xdfq0KCbK', 'revoked', 2, NOW() - INTERVAL '3 days', NULL, NOW() - INTERVAL '18 days');

ALTER SEQUENCE stations_id_seq RESTART WITH 9;

-- 6. Seed Branches (24 branches across 7 cities)
INSERT INTO branches (id, station_id, city, address, is_main, fuel_status, congestion_percentage, created_at) VALUES
(1, 1, 'الخرطوم', 'شارع أفريقيا (طريق المطار)، شارع 15', true, 'full', 65, NOW() - INTERVAL '26 days'),
(2, 1, 'الخرطوم', 'شارع الستين، حي الرياض', false, 'full', 45, NOW() - INTERVAL '20 days'),
(3, 1, 'أم درمان', 'شارع الموردة، مقابل النادي الثقافي', false, 'low', 80, NOW() - INTERVAL '18 days'),
(4, 1, 'بحري', 'شارع المعونة، وسط بحري', false, 'empty', 15, NOW() - INTERVAL '15 days'),
(5, 1, 'ود مدني', 'شارع كورنيش النيل، بالقرب من الاستاد', false, 'full', 35, NOW() - INTERVAL '12 days'),
(6, 2, 'الخرطوم', 'المنشية، ناصية شارع النيل', true, 'full', 70, NOW() - INTERVAL '25 days'),
(7, 2, 'الخرطوم', 'شارع الغابة، المنطقة الصناعية الثقيلة', false, 'low', 85, NOW() - INTERVAL '19 days'),
(8, 2, 'أم درمان', 'شارع الشنقيطي، الثورة الحارة الخامسة', false, 'full', 50, NOW() - INTERVAL '16 days'),
(9, 2, 'بحري', 'الصافية، شارع الإنقاذ', false, 'full', 40, NOW() - INTERVAL '14 days'),
(10, 2, 'بورتسودان', 'طريق ميناء البحر الأحمر، الرصيف الجنوبي', false, 'low', 60, NOW() - INTERVAL '10 days'),
(11, 3, 'بحري', 'منطقة كوبر الصناعية، شارع 1', true, 'full', 55, NOW() - INTERVAL '24 days'),
(12, 3, 'الخرطوم', 'حي الرياض، شارع المشتل', false, 'full', 62, NOW() - INTERVAL '17 days'),
(13, 3, 'أم درمان', 'شارع الأربعين، تقاطع 4', false, 'empty', 20, NOW() - INTERVAL '13 days'),
(14, 3, 'بورتسودان', 'الكورنيش الساحلي، حي الثورة', false, 'full', 48, NOW() - INTERVAL '9 days'),
(15, 4, 'أم درمان', 'تقاطع الطريق الرئيسي، سوق ليبيا', true, 'full', 75, NOW() - INTERVAL '23 days'),
(16, 4, 'أم درمان', 'شارع الوادي، الثورة الحارة 21', false, 'low', 68, NOW() - INTERVAL '15 days'),
(17, 4, 'الخرطوم', 'الصحافة زلط، المحور التجاري الرئيسي', false, 'full', 52, NOW() - INTERVAL '11 days'),
(18, 4, 'كسلا', 'الطريق الدائري، وسط مدينة كسلا', false, 'empty', 10, NOW() - INTERVAL '8 days'),
(19, 5, 'بورتسودان', 'شارع المارينا، ميناء بورتسودان', true, 'full', 58, NOW() - INTERVAL '21 days'),
(20, 5, 'بورتسودان', 'منطقة الجمارك، البوابة 7', false, 'full', 44, NOW() - INTERVAL '14 days'),
(21, 5, 'كسلا', 'محطة الطريق السريع الحدودي، الكيلو 12', false, 'low', 72, NOW() - INTERVAL '7 days'),
(22, 6, 'الأبيض', 'طريق كردفان السريع، الدائري الأبيض', true, 'empty', 0, NOW() - INTERVAL '2 days'),
(23, 7, 'ود مدني', 'تقاطع السوق الكبير، وسط المدينة', true, 'empty', 0, NOW() - INTERVAL '7 days'),
(24, 8, 'كسلا', 'ضفة نهر القاش، القطاع الشمالي', true, 'low', 25, NOW() - INTERVAL '18 days');

ALTER SEQUENCE branches_id_seq RESTART WITH 25;

-- 7. Seed Wallets (12 users + 8 stations)
INSERT INTO wallets (id, owner_type, owner_id, total_balance, allocated_balance, available_balance, created_at) VALUES
(1, 'user', 1, 0, 0, 0, NOW() - INTERVAL '30 days'),
(2, 'user', 2, 0, 0, 0, NOW() - INTERVAL '28 days'),
(3, 'user', 3, 150000, 60000, 90000, NOW() - INTERVAL '20 days'),
(4, 'user', 4, 280000, 100000, 180000, NOW() - INTERVAL '19 days'),
(5, 'user', 5, 75000, 30000, 45000, NOW() - INTERVAL '18 days'),
(6, 'user', 6, 210000, 60000, 150000, NOW() - INTERVAL '17 days'),
(7, 'user', 7, 95000, 30000, 65000, NOW() - INTERVAL '16 days'),
(8, 'user', 8, 40000, 15000, 25000, NOW() - INTERVAL '15 days'),
(9, 'user', 9, 120000, 40000, 80000, NOW() - INTERVAL '14 days'),
(10, 'user', 10, 50000, 15000, 35000, NOW() - INTERVAL '13 days'),
(11, 'user', 11, 300000, 80000, 220000, NOW() - INTERVAL '12 days'),
(12, 'user', 12, 10000, 0, 10000, NOW() - INTERVAL '11 days'),
(13, 'station', 1, 155000, 0, 155000, NOW() - INTERVAL '26 days'),
(14, 'station', 2, 100000, 0, 100000, NOW() - INTERVAL '25 days'),
(15, 'station', 3, 75000, 0, 75000, NOW() - INTERVAL '24 days'),
(16, 'station', 4, 40000, 0, 40000, NOW() - INTERVAL '23 days'),
(17, 'station', 5, 40000, 0, 40000, NOW() - INTERVAL '21 days'),
(18, 'station', 6, 0, 0, 0, NOW() - INTERVAL '2 days'),
(19, 'station', 7, 0, 0, 0, NOW() - INTERVAL '7 days'),
(20, 'station', 8, 0, 0, 0, NOW() - INTERVAL '18 days');

ALTER SEQUENCE wallets_id_seq RESTART WITH 21;

-- 8. Seed Transactions (28 transactions: 10 topups + 13 allocations + 3 station transfers + 2 withdrawals)
INSERT INTO transactions (id, reference, type, amount, from_type, from_id, from_name, to_type, to_id, to_name, branch_id, branch_name, status, created_at) VALUES
(1, 'TOP-1731140001-100001', 'topup', 150000, NULL, NULL, '', 'user', 3, 'أحمد خليل', NULL, NULL, 'completed', NOW() - INTERVAL '10 days'),
(2, 'TOP-1731140002-100002', 'topup', 280000, NULL, NULL, '', 'user', 4, 'محمد نور', NULL, NULL, 'completed', NOW() - INTERVAL '10 days' + INTERVAL '1 hour'),
(3, 'TOP-1731140003-100003', 'topup', 75000, NULL, NULL, '', 'user', 5, 'فاطمة آدم', NULL, NULL, 'completed', NOW() - INTERVAL '10 days' + INTERVAL '2 hours'),
(4, 'TOP-1731140004-100004', 'topup', 210000, NULL, NULL, '', 'user', 6, 'عمر حسن', NULL, NULL, 'completed', NOW() - INTERVAL '10 days' + INTERVAL '3 hours'),
(5, 'TOP-1731140005-100005', 'topup', 95000, NULL, NULL, '', 'user', 7, 'سارة إبراهيم', NULL, NULL, 'completed', NOW() - INTERVAL '10 days' + INTERVAL '4 hours'),
(6, 'TOP-1731140006-100006', 'topup', 40000, NULL, NULL, '', 'user', 8, 'علي عثمان', NULL, NULL, 'completed', NOW() - INTERVAL '10 days' + INTERVAL '5 hours'),
(7, 'TOP-1731140007-100007', 'topup', 120000, NULL, NULL, '', 'user', 9, 'منى صديق', NULL, NULL, 'completed', NOW() - INTERVAL '10 days' + INTERVAL '6 hours'),
(8, 'TOP-1731140008-100008', 'topup', 50000, NULL, NULL, '', 'user', 10, 'خالد إدريس', NULL, NULL, 'completed', NOW() - INTERVAL '10 days' + INTERVAL '7 hours'),
(9, 'TOP-1731140009-100009', 'topup', 300000, NULL, NULL, '', 'user', 11, 'زينب طارق', NULL, NULL, 'completed', NOW() - INTERVAL '10 days' + INTERVAL '8 hours'),
(10, 'TOP-1731140010-100010', 'topup', 10000, NULL, NULL, '', 'user', 12, 'مصطفى بكري', NULL, NULL, 'completed', NOW() - INTERVAL '10 days' + INTERVAL '9 hours'),
(11, 'TRF-1731140011-100011', 'allocation', 35000, 'user', 3, 'أحمد خليل', 'station', 1, 'شركة النيل للبترول', 1, 'الخرطوم - شارع أفريقيا (طريق المطار)، شارع 15', 'completed', NOW() - INTERVAL '9 days'),
(12, 'TRF-1731140012-100012', 'allocation', 25000, 'user', 3, 'أحمد خليل', 'station', 2, 'شركة بشائر للطاقة', NULL, NULL, 'completed', NOW() - INTERVAL '8 days'),
(13, 'TRF-1731140013-100013', 'allocation', 60000, 'user', 4, 'محمد نور', 'station', 1, 'شركة النيل للبترول', 2, 'الخرطوم - شارع الستين، حي الرياض', 'completed', NOW() - INTERVAL '8 days' + INTERVAL '3 hours'),
(14, 'TRF-1731140014-100014', 'allocation', 40000, 'user', 4, 'محمد نور', 'station', 3, 'بتروتشاينا السودان', NULL, NULL, 'completed', NOW() - INTERVAL '7 days'),
(15, 'TRF-1731140015-100015', 'allocation', 30000, 'user', 5, 'فاطمة آدم', 'station', 2, 'شركة بشائر للطاقة', 6, 'الخرطوم - المنشية، ناصية شارع النيل', 'completed', NOW() - INTERVAL '6 days'),
(16, 'TRF-1731140016-100016', 'allocation', 35000, 'user', 6, 'عمر حسن', 'station', 4, 'النهضة لخدمات الوقود', NULL, NULL, 'completed', NOW() - INTERVAL '6 days' + INTERVAL '4 hours'),
(17, 'TRF-1731140017-100017', 'allocation', 25000, 'user', 6, 'عمر حسن', 'station', 5, 'بترول البحر الأحمر', 19, 'بورتسودان - شارع المارينا، ميناء بورتسودان', 'completed', NOW() - INTERVAL '5 days'),
(18, 'TRF-1731140018-100018', 'allocation', 30000, 'user', 7, 'سارة إبراهيم', 'station', 3, 'بتروتشاينا السودان', NULL, NULL, 'completed', NOW() - INTERVAL '4 days'),
(19, 'TRF-1731140019-100019', 'allocation', 15000, 'user', 8, 'علي عثمان', 'station', 1, 'شركة النيل للبترول', 3, 'أم درمان - شارع الموردة، مقابل النادي الثقافي', 'completed', NOW() - INTERVAL '4 days' + INTERVAL '2 hours'),
(20, 'TRF-1731140020-100020', 'allocation', 40000, 'user', 9, 'منى صديق', 'station', 2, 'شركة بشائر للطاقة', NULL, NULL, 'completed', NOW() - INTERVAL '3 days'),
(21, 'TRF-1731140021-100021', 'allocation', 15000, 'user', 10, 'خالد إدريس', 'station', 5, 'بترول البحر الأحمر', NULL, NULL, 'completed', NOW() - INTERVAL '2 days'),
(22, 'TRF-1731140022-100022', 'allocation', 50000, 'user', 11, 'زينب طارق', 'station', 1, 'شركة النيل للبترول', 5, 'ود مدني - شارع كورنيش النيل، بالقرب من الاستاد', 'completed', NOW() - INTERVAL '1 day'),
(23, 'TRF-1731140023-100023', 'allocation', 30000, 'user', 11, 'زينب طارق', 'station', 2, 'شركة بشائر للطاقة', NULL, NULL, 'completed', NOW() - INTERVAL '12 hours'),
(24, 'STF-1731140024-100024', 'station_transfer', 15000, 'station', 1, 'شركة النيل للبترول', 'station', 2, 'شركة بشائر للطاقة', NULL, NULL, 'completed', NOW() - INTERVAL '5 days' + INTERVAL '5 hours'),
(25, 'STF-1731140025-100025', 'station_transfer', 10000, 'station', 2, 'شركة بشائر للطاقة', 'station', 3, 'بتروتشاينا السودان', NULL, NULL, 'completed', NOW() - INTERVAL '3 days' + INTERVAL '6 hours'),
(26, 'STF-1731140026-100026', 'station_transfer', 5000, 'station', 3, 'بتروتشاينا السودان', 'station', 4, 'النهضة لخدمات الوقود', NULL, NULL, 'completed', NOW() - INTERVAL '1 day' + INTERVAL '4 hours'),
(27, 'WTH-1731140027-100027', 'withdrawal', 25000, 'station', 1, 'شركة النيل للبترول', NULL, NULL, '', NULL, NULL, 'completed', NOW() - INTERVAL '5 days'),
(28, 'WTH-1731140028-100028', 'withdrawal', 15000, 'station', 2, 'شركة بشائر للطاقة', NULL, NULL, '', NULL, NULL, 'completed', NOW() - INTERVAL '3 days');

ALTER SEQUENCE transactions_id_seq RESTART WITH 29;

-- 9. Seed Cash Withdrawals (6 withdrawals)
INSERT INTO cash_withdrawals (id, station_id, amount, status, external_reference, reviewed_by, reviewed_at, rejection_reason, created_at) VALUES
(1, 1, 25000, 'completed', 'BOK-20261109-441820', 1, NOW() - INTERVAL '5 days', NULL, NOW() - INTERVAL '5 days' - INTERVAL '4 hours'),
(2, 2, 15000, 'completed', 'ONB-20261111-903154', 1, NOW() - INTERVAL '3 days', NULL, NOW() - INTERVAL '3 days' - INTERVAL '3 hours'),
(3, 1, 20000, 'approved', NULL, 2, NOW() - INTERVAL '1 day', NULL, NOW() - INTERVAL '1 day' - INTERVAL '5 hours'),
(4, 2, 10000, 'pending', NULL, NULL, NULL, NULL, NOW() - INTERVAL '6 hours'),
(5, 3, 15000, 'pending', NULL, NULL, NULL, NULL, NOW() - INTERVAL '2 hours'),
(6, 4, 50000, 'rejected', NULL, 1, NOW() - INTERVAL '2 days', 'المبلغ المطلوب يتجاوز الحد اليومي للسحب النقدي للمحطة', NOW() - INTERVAL '2 days' - INTERVAL '6 hours');

ALTER SEQUENCE cash_withdrawals_id_seq RESTART WITH 7;
