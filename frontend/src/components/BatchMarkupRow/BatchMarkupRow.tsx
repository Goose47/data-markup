import { block } from "../../utils/block";
import { BatchMarkupType } from "../../utils/types";
import "./BatchMarkupRow.scss";

const b = block("batch-markup-row");

export const BatchMarkupRow = ({ markup }: { markup: BatchMarkupType }) => {
  return <div className={b()}>{markup.data}</div>;
};
