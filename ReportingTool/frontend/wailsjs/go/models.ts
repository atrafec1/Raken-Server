export namespace dto {
	
	export class PayrollEntry {
	    EmployeeCode: string;
	    CurrentDate: string;
	    CraftLevel: string;
	    JobNumber: string;
	    Phase: string;
	    CostCode: string;
	    ChangeOrder: string;
	    RegularHours: number;
	    OvertimeHours: number;
	    PremiumHours: number;
	    Day: number;
	    EquipmentCode: string;
	    DownFlag: string;
	    SpecialPayType: string;
	    SpecialPayCode: string;
	    SpecialUnits: number;
	    SpecialRate: number;
	    CostCodeDivision: string;
	
	    static createFrom(source: any = {}) {
	        return new PayrollEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.EmployeeCode = source["EmployeeCode"];
	        this.CurrentDate = source["CurrentDate"];
	        this.CraftLevel = source["CraftLevel"];
	        this.JobNumber = source["JobNumber"];
	        this.Phase = source["Phase"];
	        this.CostCode = source["CostCode"];
	        this.ChangeOrder = source["ChangeOrder"];
	        this.RegularHours = source["RegularHours"];
	        this.OvertimeHours = source["OvertimeHours"];
	        this.PremiumHours = source["PremiumHours"];
	        this.Day = source["Day"];
	        this.EquipmentCode = source["EquipmentCode"];
	        this.DownFlag = source["DownFlag"];
	        this.SpecialPayType = source["SpecialPayType"];
	        this.SpecialPayCode = source["SpecialPayCode"];
	        this.SpecialUnits = source["SpecialUnits"];
	        this.SpecialRate = source["SpecialRate"];
	        this.CostCodeDivision = source["CostCodeDivision"];
	    }
	}
	export class Warning {
	    Message: string;
	    WarningType: string;
	
	    static createFrom(source: any = {}) {
	        return new Warning(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Message = source["Message"];
	        this.WarningType = source["WarningType"];
	    }
	}
	export class PayrollEntryResult {
	    Entries: PayrollEntry[];
	    Warnings: Warning[];
	
	    static createFrom(source: any = {}) {
	        return new PayrollEntryResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Entries = this.convertValues(source["Entries"], PayrollEntry);
	        this.Warnings = this.convertValues(source["Warnings"], Warning);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

