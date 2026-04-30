import type { Meta, StoryObj } from "@storybook/react";
import { CertificateStatusPill } from "../app/compliance/tds/components/CertificateStatusPill";

const meta: Meta<typeof CertificateStatusPill> = {
  title: "Compliance/TDS/Certificates/CertificateStatusPill",
  component: CertificateStatusPill,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof CertificateStatusPill>;

export const Generated: Story = { args: { status: "GENERATED" } };
export const Pending: Story = { args: { status: "PENDING" } };
export const Issued: Story = { args: { status: "ISSUED" } };
export const Revoked: Story = { args: { status: "REVOKED" } };
