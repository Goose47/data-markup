import { useState } from "react";
import { block } from "../../utils/block";
import "./BatchCreate.scss";
import { BatchForm } from "../../components/BatchForm/BatchForm";
import { CircleInfoFill } from "@gravity-ui/icons";
import { createBatch, linkBatchToMarkupType } from "../../utils/requests";
import { toaster } from "@gravity-ui/uikit/toaster-singleton";

const b = block("batch-create");

export const BatchCreate = () => {
  const [name, setName] = useState<string>("");

  const [batchType, setBatchType] = useState("");
  const [overlaps, setOverlaps] = useState("1");
  const [priority, setPriority] = useState("6");
  const [selectedType, setSelectedType] = useState("");
  const [file, setFile] = useState<File>();

  const handleCreate = () => {
    if (
      !file ||
      isNaN(parseInt(overlaps)) ||
      isNaN(parseInt(priority)) ||
      isNaN(parseInt(batchType)) ||
      isNaN(parseInt(selectedType))
    ) {
      toaster.add({
        title: "Произошла ошибка",
        name: "Все поля формы обязательны к заполнению",
        content: "Все поля формы обязательны к заполнению",
        theme: "danger",
      });
      return;
    }
    createBatch({
      name: name,
      overlaps: parseInt(overlaps),
      priority: parseInt(priority),
      markups: file,
      type_id: parseInt(batchType),
    }).then((data) => {
      linkBatchToMarkupType(data.id, parseInt(selectedType));
    });
  };

  return (
    <div className={b()}>
      <BatchForm
        submit={handleCreate}
        title={"Добавление типа разметки"}
        states={{
          name,
          setName,
          batchType,
          setBatchType,
          overlaps,
          setOverlaps,
          priority,
          setPriority,
          selectedType,
          setSelectedType,
          file,
          setFile,
        }}
        description={
          <>
            Заполните соответствующую форму и вы сможете добавить этот тип
            разметки к новым проектам.
            <br />
            <br />
            <CircleInfoFill></CircleInfoFill> Тип разметки - это то, что будет
            видеть ассессор в качестве вариантов ответа на определенный запрос в
            проекте (batch-е). Вы в последствие сможете редактировать данный
            шаблон.
          </>
        }
        buttonText={"Добавить"}
      />
    </div>
  );
};
