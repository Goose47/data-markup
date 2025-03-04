import { useContext, useState } from "react";
import { block } from "../../utils/block";
import "./BatchCreate.scss";
import { BatchForm } from "../../components/BatchForm/BatchForm";
import { CircleInfoFill } from "@gravity-ui/icons";
import { createBatch, linkBatchToMarkupType } from "../../utils/requests";
import { toaster } from "@gravity-ui/uikit/toaster-singleton";
import { LoginContext } from "../Login/LoginContext";
import { Loader } from "@gravity-ui/uikit";
import { useNavigate } from "react-router";

const b = block("batch-create");

export const BatchCreate = () => {
  const [name, setName] = useState<string>("");

  const [batchType, setBatchType] = useState("");
  const [overlaps, setOverlaps] = useState("1");
  const [priority, setPriority] = useState("6");
  const [selectedType, setSelectedType] = useState("");
  const [file, setFile] = useState<File>();

  const navigate = useNavigate();

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
      linkBatchToMarkupType(data.id, parseInt(selectedType)).then(() => {
        toaster.add({
          title: "Успешно выполнено",
          name: "Успешно привязан новый тип разметки",
          content: "Успешно привязан новый тип разметки",
          theme: "success",
        });
        navigate("/batch/" + data.id);
      });
    });
  };

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
      <BatchForm
        submit={handleCreate}
        title={"Добавление проекта"}
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
            Заполните соответствующую форму и вы создадите новый проект.
            <br />
            <br />
            <CircleInfoFill></CircleInfoFill> Проект - это набор сущностей для оценки. Позднее вы сможете редактировать свойства проекта.
          </>
        }
        buttonText={"Добавить"}
      />
    </div>
  );
};
