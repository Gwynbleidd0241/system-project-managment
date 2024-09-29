import React, { useState, useEffect } from "react";
import MyButton from "../button/MyButton.jsx";
import styles from './RedTask.module.css';
import axios from "axios";
const RedTask = (props) => {
    const rootClasses = [styles.redTask]
    if (props.visible){
      rootClasses.push(styles.active);
    }

    const [title, setTitle] = useState(props.task && props.task.title);
    const [body, setBody] = useState(props.task && props.task.body);
    
    useEffect(() => {
        if (props.task) {
          setTitle(props.task.title);
          setBody(props.task.body);
        } else {
          setTitle('');
          setBody('');
        }
      }, [props.task]);
    
      if (!props.task) {
        return null;
      }

      
      const handleSave = () => {
        axios.put('http://localhost:8080/api/task', {
          id: props.task.ID,
          email: localStorage.getItem('login'),
          title: title,
          body: body,
          token: localStorage.getItem('token')
        })
        .then(response => {
          const status = response.status;
          if (status) {
            props.onSaveTask({ id: props.task.ID, title, body });
            props.setVisible(false);
          } else {
            console.error('Ошибка редактирования задачи');
          }
        })
        .catch(error => {
          console.error(error);
        });
      };
    
      const handleDelete = () => {
        axios.delete('http://localhost:8080/api/task', {
          params: {
            id: props.task.ID,
          }
        })
        .then(response => {
          const status = response.data.status;
          if (status) {
            props.onDeleteTask(props.task.ID);
            props.setVisible(false);
          } else {
            console.error('Ошибка удаления задачи');
          }
        })
        .catch(error => {
          console.error(error);
        });
      };

    return (
        <div className={rootClasses.join(" ")} onClick={() => props.setVisible(false)}>
          <div className={styles.redTaskContent} onClick={(e) => e.stopPropagation()}>
            <div className={styles.redTask__header}>
              <h2>Редактировать задачу</h2>
            </div>
            <div className={styles.redTask__body}>
              <input type="text" value={title} onChange={(e) => setTitle(e.target.value)} placeholder="Название задачи" />
              <textarea value={body} onChange={(e) => setBody(e.target.value)} placeholder="Описание задачи" />
            </div>
            <div className={styles.redTask__footer}>
              <MyButton onClick={handleSave}>Сохранить</MyButton>
              <MyButton onClick={handleDelete}>Удалить</MyButton>
            </div>
          </div>
        </div>
      );
};

export default RedTask;