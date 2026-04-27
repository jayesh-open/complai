import type { Meta, StoryObj } from "@storybook/react";
import { SignedJsonViewer } from "../app/compliance/e-invoicing/components/SignedJsonViewer";

const meta: Meta<typeof SignedJsonViewer> = {
  title: "Compliance/SignedJsonViewer",
  component: SignedJsonViewer,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof SignedJsonViewer>;

const SAMPLE_JSON = JSON.stringify(
  {
    Version: "1.1",
    TranDtls: { TaxSch: "GST", SupTyp: "B2B", RegRev: "N" },
    DocDtls: { Typ: "INV", No: "INV-2026-01000", Dt: "26/04/2026" },
    SellerDtls: { Gstin: "29AABCA1234A1Z5", LglNm: "Acme Manufacturing Ltd" },
    BuyerDtls: { Gstin: "07AABCD3456D1Z9", LglNm: "Delhi Distributors Ltd" },
    ValDtls: {
      AssVal: 90000,
      IgstVal: 16200,
      CgstVal: 0,
      SgstVal: 0,
      TotInvVal: 106200,
    },
    ItemList: [
      {
        SlNo: "1",
        PrdDesc: "Laptop Computer (15.6 inch)",
        HsnCd: "84713010",
        Qty: 2,
        Unit: "NOS",
        UnitPrice: 45000,
        TotAmt: 90000,
        GstRt: 18,
      },
    ],
  },
  null,
  2
);

export const Default: Story = {
  args: { json: SAMPLE_JSON },
};
