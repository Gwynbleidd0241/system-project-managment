import React, { useState } from "react";
import cl from './newTaskModal.module.css'
import axios from 'axios';

const NewTaskModal = ({ visible, setVisible, onCreateTask }) => {
  const rootClasses = [cl.newTaskModal]
  if (visible){
    rootClasses.push(cl.active);
  }
  const [title, setTitle] = useState('');
  const [body, setBody] = useState('');

  const handleSubmit = (e) => {
    console.log(localStorage.getItem('login'))
    e.preventDefault();
    axios.post('http://localhost:8080/api/task', {
      title: title,
      body: body,
      email: localStorage.getItem('login'),
      token: localStorage.getItem('token'),
    })
    .then(response => {
      const userId = response.data.id;
      if (!userId) {
        console.error('Ошибка создания задачи');
      } else {
        const newTask = {
          ID: userId,
          title: title,
          body: body,
        };
        onCreateTask(newTask);
        setTitle('');
        setBody('');
        setVisible(false);
      }
    })
    .catch(error => {
      console.error(error);
    });
  };

  return (
    <div className={rootClasses.join(" ")} onClick={() => setVisible(false)}>
      <div className={cl.newTaskModalContent} onClick={(e) => e.stopPropagation()}>
        <h2>Создать новую задачу</h2>
        <form onSubmit={handleSubmit}>
          <label>
            Название задачи:
            <input type="text" value={title} onChange={(e) => setTitle(e.target.value)} />
          </label>
          <label>
            Описание задачи:
            <textarea value={body} onChange={(e) => setBody(e.target.value)} />
          </label>
          <button type="submit">Создать задачу</button>
        </form>
      </div>
    </div>
  );
};

export default NewTaskModal;