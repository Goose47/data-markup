import { useContext, useEffect, useState } from "react";
import { block } from "../../utils/block";
import "./MyBatchesCards.scss";
import { CircleInfoFill } from "@gravity-ui/icons";
import { BatchCard } from "../../components/BatchCard/BatchCard";
import { BatchCardType } from "../../utils/types";
import { getAvailableBatches } from "../../utils/requests";
import { LoginContext } from "../Login/LoginContext";
import { Loader } from "@gravity-ui/uikit";

const b = block("my-batches-cards");

const _ = require("lodash");

export const MyBatchesCards = () => {
  const [batches, setBatches] = useState<BatchCardType[]>([]);

  const handleUpdateBatch = (index: number, batch: BatchCardType) => {
    const batchesCopy = _.cloneDeep(batches);
    batchesCopy[index] = batch;
    setBatches(batchesCopy);
  };

  const [rerenderState, setRerenderState] = useState(1);
  useEffect(() => {
    getAvailableBatches().then((value: BatchCardType[]) => {
      setBatches(value);
    });
  }, [rerenderState]);

  const loginContext = useContext(LoginContext);
  if (loginContext.loading) {
    return (
      <div
        style={{
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
          height: 500,
        }}
      >
        <Loader></Loader>
      </div>
    );
  }
  if (loginContext.userRole !== "admin") {
    return (
      <div className={b()}>
        <h1>Отказано в доступе</h1>
      </div>
    );
  }

  return (
    <div className={b()}>
      <h1>Существующие проекты (batch'и)</h1>
      <p>
        <CircleInfoFill></CircleInfoFill> Здесь отображены те проекты, к которым
        у вас есть доступ в системе. Для того, чтобы проект был в работе у
        ассессора, убедитесь, что он находится в активном статусе.
      </p>

      <div className={b("list")}>
        {batches?.map((batch, index) => (
          <BatchCard
            batch={batch}
            handleUpdateBatch={(batch: BatchCardType) => {
              handleUpdateBatch(index, batch);
            }}
            triggerRerender={() => setRerenderState((v) => v + 1)}
          />
        ))}
      </div>
    </div>
  );
};
