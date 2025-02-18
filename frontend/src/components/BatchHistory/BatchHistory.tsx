import { useEffect, useState } from "react";
import { block } from "../../utils/block";
import "./BatchHistory.scss";
import { MarkupType } from "../../utils/types";
import { getAvailableMarkupTypes } from "../../utils/requests";
import { BatchHistoryCard } from "../BatchHistoryCard/BatchHistoryCard";
import { ArrowDown } from "@gravity-ui/icons";
const b = block("batch-history");

export const BatchHistory = ({ batchId }: { batchId: number }) => {
  const [history, setHistory] = useState<MarkupType[]>([]);

  useEffect(() => {
    getAvailableMarkupTypes(batchId).then((data: MarkupType[]) => {
      setHistory(data);
    });
  }, [batchId]);

  return (
    <div className={b()}>
      {history.map((historyItem, index) => (
        <>
          <BatchHistoryCard item={historyItem} />
          {index + 1 !== history.length && (
            <div className={b("arrow")}>
              <ArrowDown width={30} height={30}></ArrowDown>
            </div>
          )}
        </>
      ))}
    </div>
  );
};
