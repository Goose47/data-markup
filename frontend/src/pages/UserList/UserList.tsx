import { useEffect, useMemo, useState } from "react";
import { block } from "../../utils/block";
import "./UserList.scss";
import { UserListType } from "../../utils/types";
import { getAllProfiles } from "../../utils/requests";
import { Table } from "@gravity-ui/uikit";

const b = block("user-list");

export const UserList = () => {
  const [users, setUsers] = useState<UserListType[]>([]);

  useEffect(() => {
    getAllProfiles().then((data: UserListType[]) => {
      setUsers(
        data.map((el) => {
          return {
            ...el,
            created_at: new Date(el.created_at).toLocaleString("ru-RU"),
            precision:
              el.assessment_count + el.assessment_count2 !== 0
                ? (
                    (100 *
                      (el.correct_assessment_count +
                        el.correct_assessment_count2)) /
                    (el.assessment_count + el.assessment_count2)
                  ).toFixed(2) + "%"
                : "-",
          };
        })
      );
    });
  }, []);

  const columns = useMemo(() => {
    return [
      { id: "email", name: "E-mail" },
      { id: "precision", name: "Точность ответов" },
      {
        id: "assessment_count",
        name: "Ответов по поиску",
      },
      { id: "correct_assessment_count", name: "Из них верных" },
      { id: "assessment_count2", name: "Ответов по сравнению" },
      { id: "correct_assessment_count2", name: "Из них верных" },
      { id: "created_at", name: "Зарегистрирован" },
    ];
  }, []);

  return (
    <div className={b()}>
      <Table columns={columns} data={users}></Table>
    </div>
  );
};
