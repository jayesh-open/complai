import type { Meta, StoryObj } from "@storybook/react";
import { useState } from "react";
import { FilingConfirmationModal } from "@complai/ui-components";

const meta: Meta<typeof FilingConfirmationModal> = {
  title: "Compliance/FilingConfirmationModal",
  component: FilingConfirmationModal,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof FilingConfirmationModal>;

export const Default: Story = {
  render: () => {
    const [sign, setSign] = useState<"dsc" | "evc" | null>(null);
    return (
      <FilingConfirmationModal
        open={true} onClose={() => {}} onConfirm={() => alert("Filed!")}
        title="File GSTR-3B for April 2026"
        details={[
          { label: "Period", value: "April 2026 (FY 2026-27)" },
          { label: "GSTIN", value: "29AABCA1234A1Z5 (Karnataka)" },
          { label: "Tax payable", value: "₹12,45,678" },
        ]}
        warningText="This action is irreversible. Once filed, you cannot revise this return."
        signMethod={sign} onSignMethodChange={setSign}
      />
    );
  },
};
