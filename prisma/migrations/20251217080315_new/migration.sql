/*
  Warnings:

  - You are about to drop the `products` table. If the table is not empty, all the data it contains will be lost.

*/
-- DropForeignKey
ALTER TABLE "products" DROP CONSTRAINT "products_created_by_fkey";

-- DropForeignKey
ALTER TABLE "products" DROP CONSTRAINT "products_deleted_by_fkey";

-- DropForeignKey
ALTER TABLE "products" DROP CONSTRAINT "products_updated_by_fkey";

-- DropForeignKey
ALTER TABLE "products" DROP CONSTRAINT "products_user_id_fkey";

-- DropTable
DROP TABLE "products";
