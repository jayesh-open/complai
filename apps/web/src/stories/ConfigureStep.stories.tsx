import type { Meta, StoryObj } from "@storybook/react";
import { useState } from "react";
import { ConfigureStep } from "../app/compliance/itr/bulk/new/components/ConfigureStep";

const meta: Meta<typeof ConfigureStep> = {
  title: "Compliance/ITR/BulkWizard/ConfigureStep",
  component: ConfigureStep,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof ConfigureStep>;

export const Default: Story = {
  render: () => {
    const [batchName, setBatchName] = useState("");
    const [autoFetchAIS, setAutoFetchAIS] = useState(true);
    const [sendMagicLinks, setSendMagicLinks] = useState(true);
    return (
      <ConfigureStep
        batchName={batchName}
        setBatchName={setBatchName}
        autoFetchAIS={autoFetchAIS}
        setAutoFetchAIS={setAutoFetchAIS}
        sendMagicLinks={sendMagicLinks}
        setSendMagicLinks={setSendMagicLinks}
        selectedCount={5}
        onBack={() => {}}
        onNext={() => {}}
      />
    );
  },
};

export const Prefilled: Story = {
  render: () => (
    <ConfigureStep
      batchName="Tax Year 2026-27 — Engineering"
      setBatchName={() => {}}
      autoFetchAIS={true}
      setAutoFetchAIS={() => {}}
      sendMagicLinks={false}
      setSendMagicLinks={() => {}}
      selectedCount={12}
      onBack={() => {}}
      onNext={() => {}}
    />
  ),
};
