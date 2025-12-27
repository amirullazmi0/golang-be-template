-- AlterTable
ALTER TABLE "users" ADD COLUMN     "verification_expiry" TIMESTAMP(3),
ADD COLUMN     "verification_token" TEXT,
ALTER COLUMN "is_active" SET DEFAULT false;
