// model/task.model.ts

export interface Task {
    id: string;
    title: string;
    description: string;
    startDate: string;
    endDate: string;
    assignedTo: string;
    status: string;
    projectId: string; // Dodaj ID projekta kome zadatak pripada
  }
  