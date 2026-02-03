-- Enable Case Insensitive Text Support
CREATE EXTENSION IF NOT EXISTS "citext";

-- Enable UUID Support: uncomment to use uuid-ossp instead of pgcrypto
-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- SELECT uuid_generate_v4();