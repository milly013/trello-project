import { Component, OnInit } from '@angular/core';
import { Project, ProjectService } from '../service/project.service';
import { nextTick } from 'process';
import { CommonModule } from '@angular/common';
import { RouterLink, RouterOutlet } from '@angular/router';


@Component({
  selector: 'app-project-list',
  standalone: true,
  templateUrl: './project-list.component.html',
  styleUrls: ['./project-list.component.css'],
  imports: [CommonModule,RouterLink,RouterOutlet,]
})
export class ProjectListComponent implements OnInit {
  projects: Project[] = [];

  constructor(private projectService: ProjectService) {}

  ngOnInit(): void {
    this.projectService.getProjects().subscribe({
      next: (data: Project[]) => {
        this.projects = data;
      },
      error: (error: any) => {
        console.error('Error fetching projects', error);
      },
      complete: () => {
        console.log('Fetching projects complete');
      }
    });
  }
  
}
